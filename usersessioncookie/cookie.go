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

// TODO: probably need cookie versions

func (c *CookieSessionManager) ReadSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("ReadSession called for reading cookie:" + c.id)
	cookie, err := r.Cookie(c.id) // TODO we should check this against the domain being set yesireebob
	if err != nil {
		return
	}
	defaultLogger.Debug("Cookie captured: " + cookie.Value)
	valueBytes, _ := base64.StdEncoding.DecodeString(cookie.Value)
	defaultLogger.Debug("Readsession decoded string: " + string(valueBytes))
	value := new(cookieValue)
	json.Unmarshal(valueBytes, value)
	err = c.signer.PublicKey().Verify([]byte(value.StateString), &value.Sig)
	if err != nil {
		defaultLogger.Debug(err.Error())
		return
	}
	defaultLogger.Debug("Readsession captured string: " + value.StateString)
	if userState, ok := userStateCoord.(ReqBySess); ok {
		userState.StringToState(value.StateString)
	}
}

// Factored out of UpdateSession so we can test it + generate value strings to test javascript
func (c *CookieSessionManager) generateCookieValue(stateString string) (string, error) {
	signature, _ := c.signer.Sign(rand.Reader, []byte(stateString))
	value := &cookieValue{StateString: stateString, Sig: *signature}
	valueBytes, err := json.Marshal(value)
	defaultLogger.Debug("Signed Cookie: " + string(valueBytes))
	// the whole value
	valueString := base64.StdEncoding.EncodeToString(valueBytes)
	if err != nil {
		return "", err
	}
	defaultLogger.Debug("Cookie value will be: " + valueString)
	return valueString, nil
}

func (c *CookieSessionManager) UpdateSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("UpdateSession called")
	if userState, ok := userStateCoord.(ReqBySess); ok {
		stateString, duration := userState.StateToString()
		defaultLogger.Debug("CookieManager received cookie value string form user: " + stateString)
		valueString, err := c.generateCookieValue(stateString)
		if err != nil {
			defaultLogger.Error(err.Error())
			return
		}
		t := time.Now().Add(duration)
		if duration == 0*time.Second {
			t = time.Time{}
		}
		cookie := &http.Cookie{
			Name:     c.id,
			Value:    valueString,
			Domain:   c.domain,
			SameSite: http.SameSiteStrictMode,
			Path:     c.path,
			Expires:  t,
		}
		defaultLogger.Debug("Cookie to set: " + cookie.String())
		http.SetCookie(w, cookie)
	}
}
