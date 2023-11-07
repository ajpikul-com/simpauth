// Add Validate Option TODO (can't unless we're using different public key)
export function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        const part = parts.pop();
        if (typeof (part) === "undefined")
            return ""; // TODO can we returned undefined
        const encoded = part.split(';').shift();
        if (typeof (encoded) === "undefined")
            return ""; // TODO return something better
        return encoded;
    }
    return "";
}
export function decodeCookie(cookie) {
    if (cookie === "") {
        return null;
    }
    const decoded = atob(cookie); // TODO handle error
    return JSON.parse(decoded);
}
//# sourceMappingURL=index.js.map