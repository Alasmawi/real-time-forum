import AbstractView from "./abstract-view.js";
import { PostsController } from "../controllers/posts-controller.js";
import { CategoryController } from "../controllers/category-controller.js";
import { NewPostCardView } from "./posts/new-post-card.js";
import { UserModel } from "../models/user-model.js";

export default class PostsView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Home");
    }

    async getHtml() {
        // Get current user and categories for new post card
        const user = await UserModel.fetchCurrentUser();
        const categories = await CategoryController.loadCategories();
        
        // Generate new post card HTML
        const newPostCardHTML = user ? await NewPostCardView.generateHTML(user, categories) : '';
        
        return `
        <div class="posts-container">
            <div id="posts-list">
                ${newPostCardHTML}
                <!-- Posts will be loaded by controller -->
            </div>
        </div>
        `;
    }

    async getData() {
        // Initialize posts page using the controller
        await PostsController.initializePage();
    }
}
