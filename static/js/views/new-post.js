import AbstractView from "./abstract-view.js";

export default class NewPostView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Create New Post");
        this.categories = [];
    }

    async getHtml() {
        await this.loadCategories();
        
        return `
        <div style="border: 3px solid black; margin-top: 30px; padding: 20px;">
            <h2>Create New Post</h2>
            <form id="newpost-form">
                <label for="post-content">Content:</label><br>
                <textarea id="post-content" name="content" rows="6" cols="50" placeholder="What's on your mind?" required></textarea><br><br>
                
                <label for="categories">Categories:</label><br>
                <div id="categories-container">
                    ${this.generateCategoryCheckboxes()}
                </div><br>
                
                <input type="submit" value="Create Post">
                <button type="button" id="cancel-btn">Cancel</button>
            </form>
            
            <div id="message" style="margin-top: 20px; color: green;"></div>
            <div id="error-message" style="margin-top: 20px; color: red;"></div>
        </div>
        `;
    }

    async loadCategories() {
        try {
            const response = await fetch("/protected/v1/categories");
            
            if (response.ok) {
                const data = await response.json();
                this.categories = data;
            } else {
                console.error("Failed to load categories from API, status:", response.status);
                this.categories = [];
            }
        } catch (error) {
            console.error("Error loading categories:", error);
            this.categories = [];
        }
    }

    generateCategoryCheckboxes() {
        if (this.categories.length === 0) {
            return '<p style="color: red;">Unable to load categories. Please try again later.</p>';
        }
        
        return this.categories.map(category => `
            <input type="checkbox" id="category-${category.id}" name="categories" value="${category.id}">
            <label for="category-${category.id}">${category.name}</label><br>
        `).join('');
    }

    async getData() {
        document.getElementById("newpost-form").addEventListener("submit", async (e) => {
            e.preventDefault();
            
            let formData = {
                "content": document.getElementById("post-content").value,
                "categories": Array.from(document.querySelectorAll('input[name="categories"]:checked')).map(checkbox => parseInt(checkbox.value))
            };

            console.log("Form Data:", formData);

            if (!formData.content) {
                this.showError("Content is required");
                return;
            }

            try {
                const response = await fetch("/protected/v1/newpost", {
                    method: 'POST',
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(formData),
                });

                console.log("Response:", response);

                if (!response.ok) {
                    const errorData = await response.json();
                    console.error("Error:", errorData);
                    this.showError("Failed to create post: " + (errorData.Error || "Unknown error"));
                } else {
                    console.log("Post created successfully");
                    this.showSuccess("Post created successfully!");
                    document.getElementById("newpost-form").reset();
                    
                    setTimeout(() => {
                        window.history.pushState(null, null, "/posts");
                        window.dispatchEvent(new PopStateEvent('popstate'));
                    }, 2000);
                }
            } catch (error) {
                console.error("Error:", error);
                this.showError("An error occurred while creating the post.");
            }
        });

        document.getElementById("cancel-btn").addEventListener("click", (e) => {
            e.preventDefault();
            window.history.pushState(null, null, "/posts");
            window.dispatchEvent(new PopStateEvent('popstate'));
        });
    }

    showError(message) {
        const errorElement = document.getElementById("error-message");
        const messageElement = document.getElementById("message");
        
        errorElement.textContent = message;
        messageElement.textContent = "";
    }

    showSuccess(message) {
        const messageElement = document.getElementById("message");
        const errorElement = document.getElementById("error-message");
        
        messageElement.textContent = message;
        errorElement.textContent = "";
    }
}
