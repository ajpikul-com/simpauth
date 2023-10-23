package googlelogin

import (
	"net/http"

	"google.golang.org/api/idtoken"

	"github.com/ajpikul-com/sutils"
	"github.com/ajpikul-com/uwho"
)

type ReqByIdent interface {
	AcceptData(map[string]interface{}) bool
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
		r.ParseMultipartForm(4096)
		cookieCSRFValue, err := r.Cookie("g_csrf_token")
		if err != nil {
			defaultLogger.Error(err.Error())
			return
		}
		if cookieCSRFValue.Value != r.Form["g_csrf_token"][0] {
			defaultLogger.Info("Under attack? csrf tokens didn't match")
			return
		}

		payload, err := idtoken.Validate(r.Context(), r.Form["credential"][0], "")
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

func (g *GoogleLogin) DefaultLoginDiv(loginEndpoint string) string {
	return `<div id="g_id_onload"
		 data-client_id="` + g.ClientID + `"
		 data-context="signin"
		 data-ux_mode="popup"
		 data-login_uri="` + loginEndpoint + `"
		 data-auto_prompt="false">
</div>
<div class="g_id_signin"
		 data-type="standard"
		 data-shape="pill"
		 data-theme="outline"
		 data-text="signin_with"
		 data-size="medium"
		 data-locale="en-US"
		 data-logo_alignment="left">
</div>`
}
func (g *GoogleLogin) DefaultLoginPortal(loginEndpoint string) http.Handler {
	return sutils.StringHandler{`<html>
	<head>
		<script src="https://accounts.google.com/gsi/client" async></script>
	</head>
	<body>
		` + g.DefaultLoginDiv(loginEndpoint) + `
	</body>
</html>`}
}
