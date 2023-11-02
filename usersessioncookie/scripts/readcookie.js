function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}
layer1 = getCookie("auth-pikul")
obj =  JSON.parse(JSON.parse(atob(layer1)).StateString)
