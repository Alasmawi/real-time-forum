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
        const wsCookie = await getCookie("ws_token");
        if (!wsCookie) {
            console.warn("WebSocket cookie not found");
            return null;
        }

        if (!hasValidCookie(wsCookie)) {    
            console.warn("WebSocket cookie has expired or is invalid");
            return null; // Cookie expired
        }

        return wsCookie.value;

    } catch (error) {
        console.error('Error getting WebSocket token:', error);
        return null;
    }
}

// Global authentication helper for privilege-based routing and content loading
window.isAuthenticated = async function () {
    try {
        // Quick client-side check first (getSessionToken already validates)
        const sessionToken = await getSessionValue();
        if (!sessionToken) {
            return false;
        }

        // Server-side validation
        const response = await fetch('/v1/checkauth', {
            credentials: 'include'  // Include HttpOnly cookies
        });

        const data = await response.json();
        return data.authenticated;
    } catch {
        return false;
    }
};

// Debug function to help troubleshoot cookie issues
// async function debugCookies() {
//     try {
//         if (!('cookieStore' in window)) {
//             console.error('Cookie Store API not supported in this browser');
//             return;
//         }

//         console.log('=== Cookie Debug Info ===');

//         const sessionCookie = await cookieStore.get('session_token');
//         console.log('Session cookie object:', sessionCookie);

//         if (sessionCookie) {
//             console.log('Cookie details:');
//             console.log('- Name:', sessionCookie.name);
//             console.log('- Value:', sessionCookie.value);
//             console.log('- Domain:', sessionCookie.domain);
//             console.log('- Path:', sessionCookie.path);
//             console.log('- Expires:', sessionCookie.expires ? new Date(sessionCookie.expires) : 'Session cookie');
//             console.log('- HttpOnly:', sessionCookie.httpOnly);
//             console.log('- Secure:', sessionCookie.secure);
//             console.log('- SameSite:', sessionCookie.sameSite);
//             console.log('- Is Expired:', await isSessionExpired());
//         }

//         console.log('Has valid session:', await hasValidSession());
//         console.log('All cookies:', await getAllCookies());
//     } catch (error) {
//         console.error('Debug error:', error);
//     }
// }

export { getCookie, getSessionValue, hasValidCookie };
