import AbstractView from "./abstract-view.js";
import { NewPostController } from "../controllers/new-post-controller.js";

export default class NewPostView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Create New Post");
    }

    async getHtml() {
        return `
        <div class="new-post-form-container">
            <h2>Create New Post</h2>
            <form id="newpost-form">
                <label for="content">Content:</label><br>
                <textarea id="content" name="content" rows="6" cols="50" placeholder="What's on your mind?" required></textarea><br><br>
                
                <label for="categories">Categories:</label><br>
                <div id="categories-container">
                    <!-- Categories will be populated by the controller -->
                </div><br>
                
                <input type="submit" value="Create Post">
                <button type="button" id="cancel-btn">Cancel</button>
            </form>
            
            <div id="message" class="new-post-message success"></div>
            <div id="error-message" class="new-post-message error"></div>
        </div>
        `;
    }

    async getData() {
        // Let the controller handle everything
        await NewPostController.initializePage();
    }
}
