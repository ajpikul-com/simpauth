export function getCookie(name: string): string {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
      const part = parts.pop();
      if (typeof(part) === "undefined") return "";
      return part.split(';').shift() as string;
  }
  return "";
}
export function getCookieAsAny(name: string): any {
  const obj = JSON.parse(JSON.parse(atob(getCookie(name))).StateString) as any;
  return obj
}
