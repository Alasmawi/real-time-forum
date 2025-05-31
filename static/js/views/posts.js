import AbstractView from "./AbstarctView.js";

export default class PostsView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Posts");
    }

    async getHtml() {
        return `
        <div id="posts-container" style="margin-top: 30px; border: 3px solid black; padding: 20px;">
                <h2>Posts</h2>
                
            <div id="posts-list"><!-- Posts will be dynamically inserted here --></div>
    </div>
    `;
    }

    async getData() {
        try {
            // Fetch posts from the backend
            const response = await fetch('/v1/posts') // Replace with your actual API endpoint
            if (!response.ok) {
                throw new Error(`Failed to fetch posts: ${response.statusText}`);
            }
    
            const posts = await response.json();
    
            // Get the posts container
            const postsList = document.getElementById('posts-list');
    
            // Loop through the posts and create HTML for each
            posts.forEach(post => {
            const postElement = document.createElement('div'); 
                // Create HTML for the post attributes
                postElement.innerHTML = `
                    <h3>Post ID: ${post.id}</h3>
                    <p><strong>User ID:</strong> ${post.user_id}</p>
                    <p><strong>Username:</strong> ${post.username}</p>
                    <p><strong>First Name:</strong> ${post.first_name}</p>
                    <p><strong>Last Name:</strong> ${post.last_name}</p>
                    <p><strong>Likes:</strong> ${post.likes}</p>
                    <p><strong>Dislikes:</strong> ${post.dislikes}</p>
                    <p><strong>Comments:</strong> ${post.comments}</p>
                    <p><strong>Categories:</strong></p>
                    <ul>
                        ${post.categories.map(category => `<li> <br>${category.name}</li>`).join('')}
                    </ul>
                `;
    
                // Append the post to the container
                postsList.appendChild(postElement);
            });
        } catch (error) {
            console.error('Error fetching posts:', error);
            const postsList = document.getElementById('posts-list');
            postsList.innerHTML = '<p>Failed to load posts.</p>';
        }
    }
    
    

}
    
//     async getData() {
//         let formData = {
//             "email": document.getElementById("email").value,
//             "username": document.getElementById("username").value,
//             "password": document.getElementById("password").value,
//             "age": document.getElementById("age").value,
//             "sex": document.getElementById("sex").value,
//         };

//         try {
//             const response = await fetch("/api/posts", {
//                 method: 'POST',
//                 headers: {
//                     "Content-Type": "application/json",
//                 },
//                 body: JSON.stringify(formData),
//             });

//             console.log("Form Data:", formData);
//             console.log("Response:", response);

//             if (!response.ok) {
//                 const errorData = await response.json();
//                 console.error("Error:", errorData);
//             } else {
//                 const data = await response.json();
//                 console.log("Registration successful:", data);
//             }
//         } catch (error) {
//             console.error("Error:", error);
//             alert("An error occurred during registration.");
//         }
//     };
// }