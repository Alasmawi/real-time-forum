import AbstractView from "./abstract-view.js";

export default class CategoriesView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Categories Test");
        this.categories = [];
    }

    async getHtml() {
        await this.loadCategories();
        
        return `
        <div style="border: 3px solid blue; margin-top: 30px; padding: 20px;">
            <h2>Categories Test Page</h2>
            <div id="categories-list">
                ${this.displayCategories()}
            </div>
            <button id="reload-btn">Reload Categories</button>
        </div>
        `;
    }

    async loadCategories() {
        try {
            console.log("Attempting to fetch categories...");
            const response = await fetch("/protected/v1/categories");
            
            console.log("Response status:", response.status);
            console.log("Response ok:", response.ok);
            
            // Let ErrorHandler process the response
            const errorHandled = await window.ErrorHandler.handleResponse(response);
            if (errorHandled) {
                this.categories = [];
                return; // ErrorHandler took care of error (redirect, show error page, etc.)
            }
            
            const data = await response.json();
            console.log("Categories response:", data);
            console.log("Categories type:", typeof data);
            console.log("Categories is array:", Array.isArray(data));
            this.categories = data;
        } catch (error) {
            console.error("Error loading categories:", error);
            this.categories = [];
        }
    }

    displayCategories() {
        if (this.categories.length === 0) {
            return '<p style="color: red;">No categories found or failed to load.</p>';
        }
        
        return `
            <h3>Found ${this.categories.length} categories:</h3>
            <ul>
                ${this.categories.map(category => `
                    <li>ID: ${category.id}, Name: ${category.name}</li>
                `).join('')}
            </ul>
        `;
    }

    async getData() {
        document.getElementById("reload-btn").addEventListener("click", async () => {
            console.log("Reloading categories...");
            await this.loadCategories();
            document.getElementById("categories-list").innerHTML = this.displayCategories();
        });
    }
}
