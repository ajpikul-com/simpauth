package usersessioncookie

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajpikul-com/uwho"
	"golang.org/x/crypto/ssh"
)

type cookieValue struct {
	StateString string
	Sig         ssh.Signature
}

/*
	func (value *cookieValue) MarshalJSON() ([]byte, error) {
		return json.Marshal(value)
	}

	func (value *cookieValue) UnmarshalJSON(data []byte) error {
		return json.Unmarshal(data, value)
	}
*/
func (c *CookieSessionManager) ReadSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("ReadSession called")
	cookie, err := r.Cookie(c.id)
	if err != nil {
		return
	}
	defaultLogger.Debug("Cookie captured: " + cookie.Value)
	valueBytes, _ := base64.StdEncoding.DecodeString(cookie.Value)
	defaultLogger.Debug("Readsession decoded string: " + string(valueBytes))
	value := new(cookieValue)
	json.Unmarshal(valueBytes, value)
	err = c.private.PublicKey().Verify([]byte(value.StateString), &value.Sig)
	if err != nil {
		defaultLogger.Debug(err.Error())
		return
	}
	defaultLogger.Debug("Readsession captured string: " + value.StateString)
	if userState, ok := userStateCoord.(ReqBySess); ok {
		userState.StringToState(value.StateString)
	}
}

func (c *CookieSessionManager) UpdateSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("UpdateSession called")
	if userState, ok := userStateCoord.(ReqBySess); ok {
		stateString, duration := userState.StateToString()
		defaultLogger.Debug("CookieManager received cookie value string form user: " + stateString)
		signature, _ := c.private.Sign(rand.Reader, []byte(stateString))
		value := &cookieValue{StateString: stateString, Sig: *signature}
		valueBytes, err := json.Marshal(value)
		defaultLogger.Debug("Signed Cookie: " + string(valueBytes))
		// the whole value
		valueString := base64.StdEncoding.EncodeToString(valueBytes)
		defaultLogger.Debug("Cookie value will be: " + valueString)
		if err != nil {
			defaultLogger.Error(err.Error()) // Is this ok
			return
		}
		t := time.Now().Add(duration)
		if duration == 0*time.Second {
			t = time.Time{}
		}
		cookie := &http.Cookie{
			Name:    c.id,
			Value:   valueString,
			Domain:  c.domain,
			Path:    c.path,
			Expires: t,
		}
		defaultLogger.Debug("Cookie to set: " + cookie.String())
		http.SetCookie(w, cookie)
	}
}
