import AbstractView from "./abstract-view.js";

export default class NewPostView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Create New Post");
        this.categories = [];
    }

    async getHtml() {
        await this.loadCategories();
        console.log("Categories when generating HTML:", this.categories); // Debug log
        
        return `
        <div style="border: 3px solid black; margin-top: 30px; padding: 20px;">
            <h2>Create New Post</h2>
            <form id="newpost-form">
                <label for="content">Content:</label><br>
                <textarea id="content" name="content" rows="6" cols="50" placeholder="What's on your mind?" required></textarea><br><br>
                
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
            const response = await fetch("/v1/categories");
            
            if (response.ok) {
                const data = await response.json();
                console.log("Categories response:", data); // Debug log
                console.log("Categories type:", typeof data); // Debug log
                console.log("Categories is array:", Array.isArray(data)); // Debug log
                // Backend returns array of category objects directly
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
        // Ensure categories are loaded if they weren't during getHtml()
        if (this.categories.length === 0) {
            await this.loadCategories();
            // Re-render categories section if categories were loaded
            if (this.categories.length > 0) {
                const categoriesContainer = document.getElementById("categories-container");
                if (categoriesContainer) {
                    categoriesContainer.innerHTML = this.generateCategoryCheckboxes();
                }
            }
        }

        document.getElementById("newpost-form").addEventListener("submit", async (e) => {
            e.preventDefault();
            
            const content = document.getElementById("content").value.trim();
            const categoryCheckboxes = document.querySelectorAll('input[name="categories"]:checked');
            const categories = Array.from(categoryCheckboxes).map(checkbox => parseInt(checkbox.value));
            
            if (!content) {
                this.showError("Content is required");
                return;
            }
            
            // Optional: Check if categories are available
            if (this.categories.length === 0) {
                this.showError("Categories are not available. Please try again later.");
                return;
            }
            
            let formData = {
                "content": content,
                "categories": categories
            };

            try {
                const response = await fetch("/v1/newpost", {
                    method: 'POST',
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(formData),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    this.showError("Failed to create post: " + (errorData.Error || "Unknown error"));
                } else {
                    const data = await response.json();
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
