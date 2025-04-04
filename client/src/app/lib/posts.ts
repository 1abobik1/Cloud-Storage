export type Post = {
    id?: number;
    title: string;
    body: string;
  };
  
  
export async function getPosts(): Promise<Post[]> {
    try {
        const response = await fetch('https://jsonplaceholder.typicode.com/posts');
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    } catch (error) {
        console.error('Failed to fetch posts:', error);
        throw error;
    }
}

export async function getPostById(id:number ){
    try{
        const response = await fetch(`https://jsonplaceholder.typicode.com/posts/${id}`)
        if (!response.ok){
            throw new Error('Network response was not ok')
        } 
                    return response.json()
    }catch(error){
        console.error('Failde to fetch post:', error)
        throw error
    }
}

export async function createPost(post:Post){
    try{
        const response = await fetch('https://jsonplaceholder.typicode.com/posts',{
            method: 'POST',
            headers:{
                'Content_Type':'application/json',
            },
            body: JSON.stringify(post),
        });
        if(!response.ok){
            throw new Error('Network response was not ok')
        }
return response.json()
    }catch(error){
        console.log("faidle to create your post",error)
        throw error
    }
}

export async function updatePost(post:Post){
    try{
        const response = await fetch('https://jsonplaceholder.typicode.com/posts/${post.id}',{
            method: 'UPDATE',
            headers:{
                'Content_Type':'application/json',
            },
            body: JSON.stringify(post),
        });
        if(!response.ok){
            throw new Error('Network response was not ok')
        }
        return response.json()

    }catch(error){
        console.log('failed to update your post',error)
        throw error
    }
}
export async function deletePost(id:number){
    try{
        const response =await fetch(`https://jsonplaceholder.typicode.com/posts/${id}`,
        {method:'DELETE'})
        if(!response.ok){
            throw new Error('Network response was not ok')
        }
    }catch(error){
        console.log('failed to delete your post',error)
        throw error
    }
}