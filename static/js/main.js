import { chatFunctionality } from "./modules/chatWebsockets.js";

function Navigate(page) {
	let container = document.getElementById("content");

	fetch("../html/" + page + ".html")
		.then(response => response.text())
		.then(html => {
			container.innerHTML = html;

			// switch statement for page-specific logic
			switch (page) {
				case "chat":
					chatFunctionality();
					break;

				case "home":
					console.log("Home page loaded.");
					break;

				case "about":
					console.log("About page loaded.");
					break;

				default:
					console.log(`No specific logic for the "${page}" page.`);
					break;
			}
		})
		.catch(error => {
			console.error("Error fetching page: ", error);
		}
		);
}

Navigate("home");

document.querySelectorAll("nav a").forEach(link => {
	link.addEventListener("click", function (event) {
		event.preventDefault();

		let page = this.getAttribute("href").substring(1);

		Navigate(page);
	});
});