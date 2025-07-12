import ChatroomView from "../views/chatroom.js";
import LoginView from "../views/login.js";
import RegisterView from "../views/register.js";
import PostsView from "../views/posts.js";
import NewPostView from "../views/new-post.js";
import ErrorView from "../views/error.js";
import LogoutView from "../views/logout.js";
import ErrorHandler from "../utils/error-handler.js";

const pathToRegex = path => {
  // Check if path ends with $ (stop character for precise parameter matching)
  const hasStopChar = path.endsWith('$');
  const cleanPath = hasStopChar ? path.slice(0, -1) : path;
  
  // Choose capture pattern based on stop character
  // [^/]+ stops at forward slash, .+ captures everything (legacy behavior)
  const paramPattern = hasStopChar ? "([^/]+)" : "(.+)";
  
  return new RegExp("^" + cleanPath.replace(/\//g, "\\/").replace(/:\w+/g, paramPattern) + "$");
};

const getParams = match => {
  const values = match.result.slice(1);
  const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

  return Object.fromEntries(keys.map((key, i) => {
    return [key, values[i]];
  }));
};

const navigateTo = url => {
  history.pushState(null, null, url);
  router();
};

const router = async () => {
  const currentPath = location.pathname;

  // const isLoginOrRegister = ['/', '/register'].includes(currentPath);
  // const isErrorRoute = /^\/error\/[45]\d{2}$/.test(currentPath);

  // if (isLoginOrRegister) {
  //   // Authenticated users cannot access login/register pages
  //   const isAuthenticated = await window.isAuthenticated();
  //   if (isAuthenticated) {
  //     navigateTo('/posts');
  //     return;
  //   }
  // } else if (!isErrorRoute) {
  //   // Unauthenticated users can only access login/register and error pages
  //   const isAuthenticated = await window.isAuthenticated();
  //   if (!isAuthenticated) {
  //     navigateTo('/');
  //     return;
  //   }
  // }

  const routes = [
    { path: "/", view: LoginView },
    { path: "/register", view: RegisterView },
    { path: "/logout", view: LogoutView },
    { path: "/chat", view: ChatroomView },
    { path: "/posts", view: PostsView },
    { path: "/newpost", view: NewPostView },
    { path: "/error/:code$", view: ErrorView },
  ];

  // Find matching route
  const potentialMatches = routes.map(route => {
    return {
      route: route,
      result: currentPath.match(pathToRegex(route.path))
    };
  });

  let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);

  if (!match) {
    // If no match found, check for 404 header
    const response = await fetch('/v1/404');
    await ErrorHandler.handleResponse(response);
    return; 
  }

  const view = new match.route.view(getParams(match));
  document.querySelector("#content").innerHTML = await view.getHtml();
  await view.getData();
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
  // Make ErrorHandler available globally
  window.ErrorHandler = ErrorHandler;

  document.body.addEventListener("click", e => {
    if (e.target.matches("[data-link]")) {
      e.preventDefault();
      navigateTo(e.target.href);
    }
  });

  router();
});
