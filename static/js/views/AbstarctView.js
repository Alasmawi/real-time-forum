export default class AbstractView {
    constructor() {
        if (this.constructor === AbstractView) {
            throw new Error("Abstract classes can't be instantiated.");
        }
    }

    setTitle(title) {
        document.title = title;
    }

    async getHtml() {
        throw new Error("Method 'getHtml()' must be implemented.");
    }

    async getData() {
        throw new Error("Method 'getHtml()' must be implemented.");
    }
}
// This is an abstract class that defines the structure for all views in the application.
