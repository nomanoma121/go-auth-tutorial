import { useState, useEffect } from "react";
import "./App.css";

import { BACKEND_ENDPOINT } from "./constants";
import { PostForm } from "./components/PostForm";
import { PostList } from "./components/PostList";
import { LoginForm } from "./components/LoginForm";
import { RegisterForm } from "./components/RegisterForm";

function App() {
  const [posts, setPosts] = useState([]);
  const [user, setUser] = useState(null);
  const [isRegister, setIsRegister] = useState(true); // 新規登録かログインかを切り替えるためのステート

  // APIから投稿データを取得する関数
  const getPosts = async () => {
    // APIからデータを取得
    const res = await fetch(`${BACKEND_ENDPOINT}/api/posts`);
    // レスポンスをJSONとして解釈
    const data = await res.json();
    // postsステートを更新
    setPosts(data);
  };

  // ログインユーザーを取得する関数
  const getAuthUser = async () => {
    // ローカルストレージからトークンを取得
    const token = localStorage.getItem("token");
    if (!token) return;
    // APIからユーザー情報を取得
    const res = await fetch(`${BACKEND_ENDPOINT}/api/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (res.ok) {
      const data = await res.json();
      // もし認証に成功すればユーザーデータが返ってくるので、それをセットする
      setUser(data);
    } else {
      // 認証に失敗したら、トークンを削除する(有効期限切れや、偽のトークンが渡された場合など)
      localStorage.removeItem("token");
    }
  };

  // useEffectを使って、このコンポーネントが描画された時に実行される処理を書く
  useEffect(() => {
    // APIから投稿データを取得
    getPosts();
    // ログインユーザーを取得
    getAuthUser();
  }, []);

  return (
    <main className="app-container">
      <h1>匿名掲示板（仮アプリ）</h1>
      {/* ログイン済みの場合はPostFormを表示 */}
      {user ? (
        <>
          <p>
            ログイン中のユーザー: {user.name} ({user.email})
          </p>
          <PostForm addPost={(post) => setPosts([post, ...posts])} />
        </>
      ) : (
        <>
          {/* ログインしていない場合はLoginFormを表示 */}
          {isRegister ? (
            <RegisterForm setUser={setUser} />
          ) : (
            <LoginForm setUser={setUser} />
          )}
          <button
            onClick={() => setIsRegister(!isRegister)}
            className="register-or-login-button"
          >
            {isRegister ? "ログイン" : "新規登録"}
          </button>
        </>
      )}
      <PostList posts={posts} />
    </main>
  );
}

export default App;
