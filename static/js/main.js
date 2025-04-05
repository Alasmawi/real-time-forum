function Navigate(page) {
	let container = document.getElementById("content");
	
	fetch("../html/" + page + ".html")
		.then(response => response.text())
		.then(html => {
			container.innerHTML = html;
		})
		.catch(error => {
			console.error("Error fetching page: ", error);
		}
	);
}

Navigate("home");

document.querySelectorAll("nav a").forEach(link => {
	link.addEventListener("click", function(event) {
		event.preventDefault();
		
		let page = this.getAttribute("href").substring(1);
		
		Navigate(page);
	});
});