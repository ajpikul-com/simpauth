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
	ClientID  string
	Processor GoogleOauthProcessor
}

type GoogleOauthProcessor interface {
	ProcessOauth(*oauth2.Userinfo)
}

func New(clientID string, processor GoogleOauthProcessor) *GoogleLogin {
	return &GoogleLogin{
		ClientID:  clientID,
		Processor: processor,
	}
}

func (g *GoogleLogin) GetLoggedOutHooks() []uwho.Hook   { return nil }
func (g *GoogleLogin) GetLoggedInHooks() []uwho.Hook    { return nil }
func (g *GoogleLogin) GetAuthorizedHooks() []uwho.Hook  { return nil }
func (g *GoogleLogin) GetAboutToLoadHooks() []uwho.Hook { return nil }

func (g *GoogleLogin) VerifyCredentials(w http.ResponseWriter, r *http.Request) uwho.UserStatus {
	if r.Method == "POST" {
		r.ParseMultipartForm(4096)
		cookieCSRFValue, err := r.Cookie("g_csrf_token")
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN
		}
		if cookieCSRFValue.Value != r.Form["g_csrf_token"][0] {
			defaultLogger.Info("Under attack? csrf tokens didn't match")
			return uwho.UNKNOWN
		}
		service, err := oauth2.NewService(r.Context())
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN
		}
		userinfoService := oauth2.NewUserinfoService(service)
		userinfoGetCall := userinfoService.Get()
		userinfo, err := userinfoGetCall.Do()
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN
		}
		defaultLogger.Info(UserinfoToString(userinfo))
		g.Processor.ProcessOauth(userinfo)
		return uwho.KNOWN
	}
	return uwho.UNKNOWN
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

func UserinfoToString(info *oauth2.Userinfo) string {
	json, err := info.MarshalJSON()
	if err != nil {
		defaultLogger.Error(err.Error())
		return ""
	}
	return string(json[:])
}
