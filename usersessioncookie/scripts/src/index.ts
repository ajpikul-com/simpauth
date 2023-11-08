export type signature = {
  Format: string;
  Blob: string;
  Rest: string | null;
}
export type cookie = {
  StateString: string;
  Sig: signature;
} 

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

export function decodeCookie(cookie: string): cookie | null {
  const cookieVal: string = getCookie(cookie)
  if (cookieVal === "") {
    return null;
  }
  const decoded: string = atob(cookieVal); // TODO handle error
  return JSON.parse(decoded) as cookie;
}
