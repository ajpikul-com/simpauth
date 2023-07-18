package googlelogin

import (
	"net/http"

	"google.golang.org/api/idtoken"

	"github.com/ajpikul-com/sutils"
	"github.com/ajpikul-com/uwho"
)

type GoogleLogin struct {
	ClientID string
}

type ReqByIdent interface {
	AcceptData(map[string]interface{}) bool
}

func New(clientID string) *GoogleLogin {
	return &GoogleLogin{
		ClientID: clientID,
	}
}

func (g *GoogleLogin) GetLoggedOutHooks() []uwho.Hook   { return nil }
func (g *GoogleLogin) GetLoggedInHooks() []uwho.Hook    { return nil }
func (g *GoogleLogin) GetAuthorizedHooks() []uwho.Hook  { return nil }
func (g *GoogleLogin) GetAboutToLoadHooks() []uwho.Hook { return nil }
func (g *GoogleLogin) TestInterface(stateManager uwho.ReqByCoord) {
	if _, ok := stateManager.(ReqByIdent); !ok {
		panic("State manager doesn't satisfied required interface")
	}
}

func (g *GoogleLogin) VerifyCredentials(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "POST" {
		r.ParseMultipartForm(4096)
		cookieCSRFValue, err := r.Cookie("g_csrf_token")
		if err != nil {
			defaultLogger.Error(err.Error())
			return false
		}
		if cookieCSRFValue.Value != r.Form["g_csrf_token"][0] {
			defaultLogger.Info("Under attack? csrf tokens didn't match")
			return false
		}

		payload, err := idtoken.Validate(r.Context(), r.Form["credential"][0], "")
		if err != nil {
			defaultLogger.Error(err.Error())
			return false
		}

		userState, ok := userStateCoord.(ReqByIdent)
		if !ok {
			defaultLogger.Error("Interface assertion error")
			return false
		}

		return userState.AcceptData(payload.Claims)
	}
	return false
}

func (g *GoogleLogin) GoHome(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Info("Redirecting")
	http.Redirect(w, r, "/", 302)
}

func (g *GoogleLogin) DefaultLoginResult(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Info("Redirecting")
	defaultLogger.Info(r.URL.Path)
	defaultLogger.Info(r.Header.Get("Referer"))
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func (g *GoogleLogin) DefaultLoginPortal(loginEndpoint string) http.Handler {
	return sutils.StringHandler{`<html>
		<head>
			<script src="https://accounts.google.com/gsi/client" async></script>
		</head>
		<body>
		<div>
		Enter your email and we'll send you a link
		</div>
		<div id="g_id_onload"
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
		</div>
		<body>
		</html>`}
}
