package usersessioncookie

import (
	"net/http"

	"github.com/ajpikul-com/uwho"
)

/*
	 time code
			t, _ := time.Now().MarshalText() // looks I forgot to add duration here
			string(t[:])
			time.Now().Add(c.expiry)
	t, err := time.Parse(time.RFC3339, splitValue[1])
	if err != nil {
		c.EndSession(userStateCoord, w, r)
		return false
	}
	if c.expiry != 0 && time.Now().After(t.Add(c.expiry)) {
		c.EndSession(userStateCoord, w, r)
		expired = true
	}
*/

func (c *CookieSessionManager) EndSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   c.id,
		Value:  "",
		MaxAge: -1,
		Domain: c.domain,
		Path:   c.path,
	})
}
