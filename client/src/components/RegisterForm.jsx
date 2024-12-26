import { useState } from "react";
import { BACKEND_ENDPOINT } from "../constants";

export function RegisterForm({
  // 親コンポーネントからログインユーザー情報をセットするための関数を受け取る
  setUser,
}) {
  const [error, setError] = useState(null);

  // フォームの送信ボタンが押された時の処理
  const handleSubmit = (e) => {
    // ページ遷移を防ぐ（デフォルトでは、フォーム送信ボタンを押すとページが遷移してしまう）
    e.preventDefault();
    // フォームの内容を取得
    const name = e.target.elements.name.value;
    const email = e.target.elements.email.value;
    const password = e.target.elements.password.value;
    // 投稿を作成
    register(name, email, password);
    // フォームを空にする
    e.target.reset();
  };

  // 投稿を作成する関数
  const register = async (name, email, password) => {
    // APIに送るデータを作成
    const payload = {
      email: email,
      password: password,
      name: name,
    };
    // APIにデータを送信
    const res = await fetch(`${BACKEND_ENDPOINT}/register`, {
      // POSTメソッドで送信
      method: "POST",
      headers: {
        // JSON形式でデータを送ることを指定
        "Content-Type": "application/json",
      },
      // JSON.stringifyでJSON文字列に変換
      // (文字列じゃないと送れないので、オブジェクトをJSON文字列に変換しています)
      body: JSON.stringify(payload),
    });
    // レスポンスをJSONとして解釈
    // ユーザー情報とトークンが返ってくる
    const data = await res.json();
    if (res.ok) {
      // ユーザー情報をセット
      setUser(data.user);
      // トークンをローカルストレージに保存
      localStorage.setItem("token", data.token);
    } else {
      // エラー情報をセット
      setError(data.message);
    }
  };

  return (
    <>
      <h2>ユーザー登録</h2>
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          name="name"
          className="auth-form__input"
          type="text"
          placeholder="ユーザー名"
        />
        <input
          name="email"
          className="auth-form__input"
          type="email"
          placeholder="メールアドレス"
        />
        <input
          name="password"
          className="auth-form__input"
          type="password"
          placeholder="パスワード"
        />
        <button type="submit" className="auth-form__submit-button">
          登録
        </button>
      </form>
      {error && <p className="auth-form__error">{error}</p>}
    </>
  );
}
