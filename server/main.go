package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "db.sqlite3"

var secret []byte
var db *sql.DB

// PostテーブルのSQLまとめ
const (
	createPostTable = `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	// 投稿の作成を行うSQL文
	insertPost = "INSERT INTO posts (content, created_at) VALUES (?, ?)"

	// 投稿の取得を行うSQL文
	selectPosts = "SELECT * FROM posts ORDER BY created_at DESC"
)

// Postは、投稿を表す構造体
type Post struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// UserテーブルのSQLまとめ
const (
	// Userテーブルの作成を行うSQL文
	createUserTable = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT UNIQUE,
			password TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	// ユーザーの作成を行うSQL文
	insertUser = "INSERT INTO users (name, email, password, created_at) VALUES (?, ?, ?, ?)"

	// ユーザーの取得を行うSQL文 (メールアドレスで取得)
	selectUserByEmail = "SELECT * FROM users WHERE email = ?"

	// ユーザーの取得を行うSQL文 (IDで取得)
	selectUserByID = "SELECT * FROM users WHERE id = ?"
)

// Userは、ユーザーを表す構造体
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// init関数は、main関数よりも先に実行される特殊な関数
func init() {
	// .envファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// secretの取得と設定
	secretRaw := os.Getenv("SECRET")
	if secretRaw == "" {
		panic("SECRETが設定されていません、環境変数を設定してください")
	}
	secret = []byte(secretRaw)

	// データベースとの接続
	db, err = sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err) // もし接続に失敗したら、プログラムを強制終了する
	}

	// データベースの接続を閉じる(init()が終了したら閉じる)
	defer db.Close()

	// Postテーブルの作成
	_, err = db.Exec(createPostTable)
	if err != nil {
		panic(err)
	}

	// Userテーブルの作成
	_, err = db.Exec(createUserTable)
	if err != nil {
		panic(err)
	}
}

// main関数は、プログラムのエントリーポイント、init()関数の実行後に実行される
func main() {
	var err error
	// データベースとの接続
	db, err = sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err) // もし接続に失敗したら、プログラムを強制終了する
	}

	// データベースの接続を閉じる(main()が終了したら閉じる)
	defer db.Close()

	// ルーティングの設定
	http.HandleFunc("/api/posts", HandleCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPosts(w, r)
		case http.MethodPost:
			HandleAuthRequire(createPost)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))


	// ルーティングの設定
	http.HandleFunc("/api/me", HandleCORS(HandleAuthRequire(getMe)))

	// ルーティングの設定
	http.HandleFunc("/login", HandleCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			login(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	// ルーティングの設定
	http.HandleFunc("/register", HandleCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			register(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	// サーバーの起動、ポート番号は8080
	fmt.Println("http://localhost:8080 でサーバーを起動します")
	http.ListenAndServe(":8080", nil)
}

// 投稿の一覧を取得する
// GET /api/posts
func getPosts(w http.ResponseWriter, r *http.Request) {
	// 投稿の取得
	rows, err := db.Query(selectPosts)
	if err != nil {
		panic(err) // もし取得に失敗したら、プログラムを強制終了する
	}
	defer rows.Close()

	// 投稿の一覧を格納する配列
	var posts = []Post{}

	// 取得した投稿を一つずつ取りだして、配列に格納する
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	// 取得した投稿をJSON形式でレスポンスする
	respondJSON(w, http.StatusOK, posts)
}

// 投稿を作成する
// POST /api/posts
func createPost(w http.ResponseWriter, r *http.Request) {
	// リクエストボディの読み込み
	var post Post
	if err := decodeBody(r, &post); err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()

	// 投稿の作成
	result, err := db.Exec(insertPost, post.Content, now)
	if err != nil {
		panic(err)
	}

	// 作成した投稿のIDを取得する
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	post.ID = int(id)
	// goのtimeでは、YYYY-MM-DD hh:mm:ssの形式でフォーマットするには、以下のようにする
	// 僕はこの書き方嫌いです！！！！！！
	post.CreatedAt = now.Format("2006-01-02 15:04:05")

	// 作成した投稿をJSON形式でレスポンスする
	respondJSON(w, http.StatusCreated, post)
}

// ユーザーを作成する
// POST /register
func register(w http.ResponseWriter, r *http.Request) {
	// リクエストボディの読み込み
	var user User
	if err := decodeBody(r, &user); err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()

	// Passwordをハッシュ化する
	// ハッシュ化操作は不可逆性があるため、一度ハッシュ化したパスワードは元に戻せない
	// そのため、ハッシュ化したパスワードをデータベースに保存すると、データベースからデータが漏洩したりしても元のパスワードを知ることができないようになる
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
		return
	}

	// ユーザーの作成
	result, err := db.Exec(insertUser, user.Name, user.Email, string(hashedPassword), now)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "そのメールアドレスは既に登録されています"})
		return
	}

	// 作成したユーザーのIDを取得する
	id, err := result.LastInsertId()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
	}
	user.ID = int(id)
	// goのtimeでは、YYYY-MM-DD hh:mm:ssの形式でフォーマットするには、以下のようにする
	user.CreatedAt = now.Format("2006-01-02 15:04:05")
	user.Password = "" // パスワードはレスポンスに含めない

	// JWTの作成
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72時間が有効期限
	}

	// ヘッダーとペイロードを結合した文字列を作成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 署名を設定する
	tokenString, err := token.SignedString(secret)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Internal Server Error"})
		return
	}

	// 作成したユーザーとJWTをJSON形式でレスポンスする
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"token": tokenString,
		"user":  user,
	})
}

// ログインする
// POST /api/login
func login(w http.ResponseWriter, r *http.Request) {
	// リクエストボディの読み込み
	var user User
	if err := decodeBody(r, &user); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "不正なリクエストです"})
		return
	}

	// ユーザーの取得
	row := db.QueryRow(selectUserByEmail, user.Email)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "メールアドレスまたはパスワードが間違っています"})
		return
	}

	// パスワードの照合
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "メールアドレスまたはパスワードが間違っています"})
		return
	}

	// JWTの作成
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72時間が有効期限
	}

	// ヘッダーとペイロードを結合した文字列を作成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 署名を設定する
	tokenString, err := token.SignedString(secret)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "不正なリクエストです"})
		return
	}

	u.Password = "" // パスワードはレスポンスに含めない

	// トークンをレスポンスする
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
		"user":  u,
	})
}

// AccessTokenからログイン中のユーザーを取得して返す
// GET /api/me
func getMe(w http.ResponseWriter, r *http.Request) {
	// コンテキストからUserIDを取得する
	userID := r.Context().Value(AuthCtxKey("user_id")).(int)

	// ユーザーの取得
	row := db.QueryRow(selectUserByID, userID)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// パスワードはレスポンスに含めない
	user.Password = ""

	// ユーザーをレスポンスする
	respondJSON(w, http.StatusOK, user)
}

// decodeBodyは、リクエストボディを構造体に変換する
// 【触るのは非推奨】
func decodeBody(r *http.Request, v interface{}) error {
	// リクエストボディの読み込み
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

// respondJSONは、JSON形式でレスポンスする
// 【触るのは非推奨】
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	// レスポンスヘッダーの設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// レスポンスボディの設定
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

// CORSを許可するミドルウェア
// 【触るのは非推奨】
func HandleCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// レスポンスヘッダーの設定
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// リクエストヘッダーの設定
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// ハンドラーの実行
		h(w, r)
	}
}

// UserIDをコンテキストに設定するためのkeyの型
// 【触るのは非推奨】
type AuthCtxKey string

// ログインが必要なハンドラー
// 【触るのは非推奨】
func HandleAuthRequire(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーからトークンを取得する
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			respondJSON(w, http.StatusUnauthorized, "no token")
			return
		}

		// Bearerトークンを取得する
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			respondJSON(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// トークン文字列を取得する
		tokenString := bearerToken[1]

		// トークンを検証する
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, err.Error())
			return
		}

		// トークンのペイロードからユーザーIDを取得する
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondJSON(w, http.StatusUnauthorized, "invalid token")
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			respondJSON(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// ユーザーIDをコンテキストに設定する
		ctx := context.WithValue(r.Context(), AuthCtxKey("user_id"), int(userID))

		// ハンドラーの実行
		h(w, r.WithContext(ctx))
	}
}
