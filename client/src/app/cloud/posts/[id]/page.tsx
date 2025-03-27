import {  getPostbyId } from "@/app/lib/posts";

export default async function PostPage({ params }: { params: { id: number } }) {
  const post = await getPostbyId(Number(params.id));

  return (
    <div>
      <h1>{post.title}</h1>
      <p>{post.body}</p>
    </div>
  );
}
