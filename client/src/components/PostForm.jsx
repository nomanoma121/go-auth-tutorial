import { BACKEND_ENDPOINT } from "../constants";
import { SendIcon } from "./SendIcon";

export function PostForm({
  // 親コンポーネントからPostを追加するための関数を受け取る
  addPost,
}) {
  // フォームの送信ボタンが押された時の処理
  const handleSubmit = (e) => {
    // ページ遷移を防ぐ（デフォルトでは、フォーム送信ボタンを押すとページが遷移してしまう）
    e.preventDefault();
    // フォームの内容を取得
    const content = e.target.elements.content.value;
    // 投稿を作成
    createPost(content);
    // フォームを空にする
    e.target.reset();
  };

  // 投稿を作成する関数
  const createPost = async (text) => {
    // APIに送るデータを作成
    const payload = {
      content: text,
    };
    // Tokenを取得
    const token = localStorage.getItem("token");
    // APIにデータを送信
    const res = await fetch(`${BACKEND_ENDPOINT}/api/posts`, {
      // POSTメソッドで送信
      method: "POST",
      headers: {
        // JSON形式でデータを送ることを指定
        "Content-Type": "application/json",
        // AuthorizationヘッダーにTokenをセット
        Authorization: `Bearer ${token}`,
      },
      // JSON.stringifyでJSON文字列に変換
      // (文字列じゃないと送れないので、オブジェクトをJSON文字列に変換しています)
      body: JSON.stringify(payload),
    });
    // レスポンスをJSONとして解釈
    // （作成した投稿データが返ってくる）
    const data = await res.json();
    if (res.ok) {
      // postsステートを更新
      addPost(data);
    }
  };

  return (
    <>
      <h2>新規投稿</h2>
      <form onSubmit={handleSubmit} className="post-form">
        <textarea name="content" rows="5" className="post-form__textarea" />
        <button type="submit" className="post-form__submit-button">
          <SendIcon />
        </button>
      </form>
    </>
  );
}
