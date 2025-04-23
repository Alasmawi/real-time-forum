// import { chatFunctionality } from "./modules/chatWebsockets.js";
// import { Router } from "./modules/routes.js/index.js";

// window.app = {}
// app.router = Router;

// window.addEventListener("DOMContentLoaded", function () {
// 	// Initialize the router
// 	app.router.init();

// 	// Add event listeners to navigation links
// 	document.querySelectorAll("nav a").forEach(link => {
// 		link.addEventListener("click", function (event) {
// 			event.preventDefault();
// 			let page = this.getAttribute("href").substring(1);
// 			Navigate(page);
// 		});
// 	});
// })

// function Navigate(page) {
// 	let container = document.getElementById("content");

// 	fetch("../html/" + page + ".html")
// 		.then(response => response.text())
// 		.then(html => {
// 			container.innerHTML = html;

// 			// switch statement for page-specific logic
// 			switch (page) {
// 				case "chat":
// 					chatFunctionality();
// 					break;

// 				case "login":
// 					break;

// 				default:
// 					console.log(`No specific logic for the "${page}" page.`);
// 					break;
// 			}
// 		})
// 		.catch(error => {
// 			console.error("Error fetching page: ", error);
// 		}
// 		);
// }

// Navigate("home");

// document.querySelectorAll("nav a").forEach(link => {
// 	link.addEventListener("click", function (event) {
// 		event.preventDefault();

// 		let page = this.getAttribute("href").substring(1);

// 		Navigate(page);
// 	});
// });