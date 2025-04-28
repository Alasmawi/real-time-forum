import ChatroomView from "../views/chatroom.js";
import LoginView from "../views/login.js";
import RegisterView from "../views/register.js";

const navigateTo = url => {
  history.pushState(null, null, url);
  router();
}

const router = async () => {
  const routes = [
    { path: "/", view: console.log("root") },
    { path: "/chat", view: ChatroomView },
    { path: "/login", view: LoginView },
    { path: "/register", view: RegisterView },
    { path: "/error", view: console.log("error") },
    
  ];

  const potentialMatches = routes.map(route => {
    return {
      route: route,
      isMatch: location.pathname === route.path || location.pathname.startsWith(route.path + "/")
    };
  })

  let match = potentialMatches.find(potentialMatch => potentialMatch.isMatch);

  if (!match) {
    match = {
      route: routes[0],
      isMatch: true
    };
  }

  const view = new match.route.view(getParams(match));
  
  document.querySelector("#content").innerHTML = await view.getHtml();

};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
  document.body.addEventListener("click", e => {
    if (e.target.matches("[data-link]")) {
      e.preventDefault();
      navigateTo(e.target.href);
    }
  });

  router();
});

export { router };