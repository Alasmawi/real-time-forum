// import ChatroomView from "../views/chatroom.js";
import LoginView from "../views/login-register/login.js";
import RegisterView from "../views/login-register/register.js";
import PostsView from "../views/home-view.js";
// import NewPostView from "../views/new-post.js"; // Removed - inline post creation now
import ErrorView from "../views/error.js";
import ErrorHandler from "./utils/errors/error-handler.js";
import { DOMControl } from "./utils/dom/layout-control.js";
import { AppInit } from "./utils/app-init.js";
import { isAuthenticated } from "./utils/session.js";

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
  if (location.pathname !== url) {
    history.pushState(null, null, url);
  }
  router();
};

const router = async () => {
  const currentPath = location.pathname;

  const routes = [
    { path: "/", view: LoginView, public: true },
    { path: "/register", view: RegisterView, public: true },
    { path: "/home", view: PostsView, public: false },
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
    const errorHandled = await ErrorHandler.handleResponse(response);
    if (errorHandled) {
      return;
    }
  }

  // Check authentication based on route's public property
  const isAuth = await isAuthenticated();
  const isErrorRoute = match.route.path.includes('/error/');
  
  if (match.route.public && isAuth && !isErrorRoute) {
    // Public route but user is authenticated (except error pages) - redirect to home
    history.pushState(null, null, '/home');
    location.reload();
    return;
  } else if (!match.route.public && !isAuth && !isErrorRoute) {
    // Private route but user not authenticated (except error pages) - redirect to login
    history.pushState(null, null, '/');
    location.reload();
    return;
  }

  // Set DOM layout mode based on route's public property
  const layoutMode = match.route.public ? 'guest' : 'authenticated';
  DOMControl.setLayoutMode(layoutMode);

  const view = new match.route.view(getParams(match));
  document.querySelector("#content").innerHTML = await view.getHtml();
  await view.getData();

  if (layoutMode === 'authenticated') {
    AppInit.updateActiveNavigation();
  }
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
  // Make ErrorHandler available globally
  window.ErrorHandler = ErrorHandler;

  document.body.addEventListener("click", e => {
    if (e.target.matches("[data-link]")) {
      e.preventDefault();
      const url = new URL(e.target.href);
      navigateTo(url.pathname);
    }
  });

  router();
});
