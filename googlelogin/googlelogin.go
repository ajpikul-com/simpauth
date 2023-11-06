package googlelogin

import (
	"io/ioutil"
	"net/http"

	"google.golang.org/api/idtoken"

	"github.com/ajpikul-com/sutils"
	"github.com/ajpikul-com/uwho"
)

type ReqByIdent interface {
	AcceptData(map[string]interface{})
}

type GoogleLogin struct {
	ClientID string
}

func New(clientID string) *GoogleLogin {
	return &GoogleLogin{
		ClientID: clientID,
	}
}

func (g *GoogleLogin) TestInterface(stateManager uwho.ReqByCoord) {
	if _, ok := stateManager.(ReqByIdent); !ok {
		panic("State manager doesn't satisfied required interface")
	}
}

func (g *GoogleLogin) VerifyCredentials(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ct := r.Header.Get("Content-Type")
		token := ""
		if r.URL.Query().Get("source") != "xhr" {
			r.ParseMultipartForm(4096)
			cookieCSRFValue, err := r.Cookie("g_csrf_token") // Don't need with xhr I don't think?
			if err != nil {
				defaultLogger.Error(err.Error())
				return
			}
			if cookieCSRFValue.Value != r.Form["g_csrf_token"][0] {
				defaultLogger.Info("Under attack? csrf tokens didn't match")
				return
			}
		}
		if ct == "application/x-www-form-urlencoded" {
			token = r.Form["credential"][0]
		} else if ct == "text/plain" {
			temp, err := ioutil.ReadAll(r.Body)
			if err != nil {
				defaultLogger.Error(err.Error())
				return
			}
			token = string(temp)
		}
		payload, err := idtoken.Validate(r.Context(), token, "")
		if err != nil {
			defaultLogger.Error(err.Error())
			return
		}
		userState, ok := userStateCoord.(ReqByIdent)
		if !ok {
			defaultLogger.Error("Interface assertion error")
			return
		}
		userState.AcceptData(payload.Claims)
	}
}

func DefaultLoginDiv(loginEndpoint string, clientID string) string {
	return `<div id="g_id_onload"
	data-client_id="` + clientID + `"
	data-context="signin"
	data-ux_mode="popup"
	data-login_uri="` + loginEndpoint + `"
	data-auto_prompt="false"
</div>

<div class="g_id_signin"
	data-type="icon"
	data-shape="circle"
	data-theme="outline"
	data-text="continue_with"
	data-size="large"
</div>`
}

func DefaultLoginPortal(loginEndpoint string, clientID string) http.Handler {
	return sutils.StringHandler{`<!DOCTYPE html>
<html>
	<head>
		<script src="https://accounts.google.com/gsi/client" async></script>
	</head>
	<body>
		` + DefaultLoginDiv(loginEndpoint, clientID) + `
	</body>
</html>`}
}
