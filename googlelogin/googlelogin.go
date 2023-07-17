package googlelogin

import (
	"net/http"

	"google.golang.org/api/oauth2/v2"

	"github.com/ajpikul-com/ilog"
	"github.com/ajpikul-com/sutils"
	"github.com/ajpikul-com/uwho"
)

var defaultLogger ilog.LoggerInterface

func init() {
	if defaultLogger == nil {
		defaultLogger = new(ilog.EmptyLogger)
	}
}

func SetDefaultLogger(newLogger ilog.LoggerInterface) {
	defaultLogger = newLogger
	defaultLogger.Info("Default Logger Set")
}

type GoogleLogin struct {
	ClientID string
}

func (g *GoogleLogin) GetLoggedOutHooks() []*uwho.Hook   { return nil }
func (g *GoogleLogin) GetLoggedInHooks() []*uwho.Hook    { return nil }
func (g *GoogleLogin) GetAuthorizedHooks() []*uwho.Hook  { return nil }
func (g *GoogleLogin) GetaboutToLoadHooks() []*uwho.Hook { return nil }

func (g *GoogleLogin) VerifyCredentials(w http.ResponseWriter, r *http.Request) (uwho.UserStatus, interface{}) {
	if r.Method == "POST" {
		r.ParseMultipartForm(4096)
		cookieCSRFValue, err := r.Cookie("g_csrf_token")
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN, nil
		}
		if cookieCSRFValue.Value != r.Form["g_csrf_token"][0] {
			defaultLogger.Info("Under attack? csrf tokens didn't match")
			return uwho.UNKNOWN, nil // BAD! Not Verified Should probably block this user
		}
		service, err := oauth2.NewService(r.Context())
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN, nil // Not validated!
		}
		userinfoService := oauth2.NewUserinfoService(service)
		userinfoGetCall := userinfoService.Get()
		userinfo, err := userinfoGetCall.Do()
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN, nil // Not Validated
		}
		return uwho.KNOWN, userinfo
	}
	return uwho.UNKNOWN, nil
}
func (g *GoogleLogin) DefaultLoginResult(w http.ResponseWriter, r *http.Request) {
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

// TODO write helper for "AboutToLoad"
