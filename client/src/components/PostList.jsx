export function PostList({
  // 親コンポーネントから投稿データを受け取る
  posts,
}) {
  return (
    <>
      <h2>投稿一覧</h2>
      <div className="post-list">
        {posts.map((post) => (
          <div key={post.id} className="post-list__item">
            <span className="post-list__item__content">{post.content}</span>
            <span className="post-list__item__date">
              {new Date(post.created_at).toLocaleString("ja-JP")}
            </span>
          </div>
        ))}
      </div>
    </>
  );
}
