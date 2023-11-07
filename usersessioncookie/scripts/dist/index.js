/* type {
  StateString: string
  Sig: {
  } */
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
    const decoded = atob(cookie); // TODO handle error
    return JSON.parse(decoded); // TODO lets see what we get here
}
export function getCookieAsAny(name) {
    const obj = JSON.parse(JSON.parse(atob(getCookie(name))).StateString);
    return obj;
}
//# sourceMappingURL=index.js.map