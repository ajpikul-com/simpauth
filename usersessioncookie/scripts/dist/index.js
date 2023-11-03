"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getCookieAsAny = exports.getCookie = void 0;
function getCookie(name) {
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
exports.getCookie = getCookie;
function getCookieAsAny(name) {
    const obj = JSON.parse(JSON.parse(atob(getCookie(name))).StateString);
    return obj;
}
exports.getCookieAsAny = getCookieAsAny;
//# sourceMappingURL=index.js.map