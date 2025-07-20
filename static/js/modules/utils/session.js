import ErrorHandler from './errors/error-handler.js';

// Get a specific cookie using Cookie Store API
async function getCookie(name) {
    if ("cookieStore" in window) {
        console.log("Browser supports cookieStore API");
    } else {
        alert("Cookie Store API not supported in this browser");
        return null;
    }
    try {
        const cookie = await cookieStore.get(name);
        return cookie || null;
    } catch (error) {
        console.error('Error getting cookie:', error);
        return null;
    }
}

// Check if session cookie has expired using actual expiry metadata
function hasValidCookie(cookie) {
    if (!cookie.expires) {
        console.log("Cookie ", cookie.name, " does not have an expiry date.");
        return true;
    }

    const expiryDate = new Date(cookie.expires).toLocaleString();

    switch (true) {
        case (Date.now() > cookie.expires):
            console.warn("Cookie ", cookie.name, " has expired");
            return false; // Cookie is expired
        case (cookie.value === ""):
            console.warn("Cookie ", cookie.name, " is empty");
            return false; // Empty cookie value is invalid
        default:
            console.log("Cookie ", cookie.name, "is valid until ", expiryDate);
            return true; // Cookie is valid
    }
}
// Get WebSocket token with implicit expiry validation
// Returns null if token is missing or expired
async function getSessionValue() {
    try {
        const sessionCookie = await getCookie("session_token");
        if (!sessionCookie) {
            console.warn("Session cookie not found");
            return null;
        }

        if (!hasValidCookie(sessionCookie)) {    
            console.warn("Session cookie has expired or is invalid");
            return null; // Cookie expired
        }

        return sessionCookie.value;

    } catch (error) {
        console.error('Error getting WebSocket token:', error);
        return null;
    }
}

// Global authentication helper for privilege-based routing and content loading
async function isAuthenticated() {
    try {
        // Client-side check only - backend validates sessions on protected routes
        const sessionToken = await getSessionValue();
        return sessionToken !== null;
    } catch (error) {
        return false;
    }
};

export { getCookie, getSessionValue, hasValidCookie, isAuthenticated };
