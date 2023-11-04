export function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        const part = parts.pop();
        if (typeof (part) === "undefined")
            return "";
        return part.split(';').shift();
    }
    return "";
}
export function getCookieAsAny(name) {
    const obj = JSON.parse(JSON.parse(atob(getCookie(name))).StateString);
    return obj;
}
//# sourceMappingURL=index.js.map