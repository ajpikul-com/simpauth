/* type {
  StateString: string
  Sig: {
  } */

// Add Validate Option TODO (can't unless we're using different public key)

export function getCookie(name: string): string { 
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    const part = parts.pop();
    if (typeof(part) === "undefined") return "";  // TODO can we returned undefined
    const encoded: string | undefined = part.split(';').shift();
    if (typeof(encoded) === "undefined") return ""; // TODO return something better
    return encoded;
  }
  return "";
}

export function decodeCookie(cookie: string): any {
    const decoded: string = atob(cookie); // TODO handle error
    return JSON.parse(decoded); // TODO lets see what we get here
}

export function getCookieAsAny(name: string): any {
  const obj = JSON.parse(JSON.parse(atob(getCookie(name))).StateString) as any;
  return obj
}
