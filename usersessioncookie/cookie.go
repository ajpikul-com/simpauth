package usersessioncookie

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajpikul-com/uwho"
	"golang.org/x/crypto/ssh"
)

type cookieValue struct {
	stateString string
	sig         ssh.Signature
}

func (value *cookieValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value)
}

func (value *cookieValue) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, value)
}

func (c *CookieSessionManager) ReadSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(c.id)
	if err != nil {
		return
	}
	value := &cookieValue{}
	value.UnmarshalJSON([]byte(cookie.Value))
	err = c.private.PublicKey().Verify([]byte(value.stateString), &value.sig)
	if err != nil {
		defaultLogger.Debug(err.Error())
		return
	}
	defaultLogger.Debug("Readsession captured string: " + value.stateString)
	if userState, ok := userStateCoord.(ReqBySess); ok {
		userState.StringToState(value.stateString)
	}
}

func (c *CookieSessionManager) UpdateSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	if userState, ok := userStateCoord.(ReqBySess); ok {
		stateString, duration := userState.StateToString()
		signature, _ := c.private.Sign(rand.Reader, []byte(stateString))
		value := &cookieValue{stateString: stateString, sig: *signature}
		valueBytes, err := value.MarshalJSON()
		if err != nil {
			defaultLogger.Error(err.Error()) // Is this ok
			return
		}
		t := time.Now().Add(duration)
		if duration == 0*time.Second {
			t = time.Time{}
		}
		http.SetCookie(w, &http.Cookie{
			Name:    c.id,
			Value:   string(valueBytes),
			Domain:  c.domain,
			Path:    c.path,
			Expires: t,
		})
	}
}
