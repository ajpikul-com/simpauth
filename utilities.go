package uwho

import "net/http"

type RedirectHome string

func (d *RedirectHome) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("Redirecting")
	http.Redirect(w, r, string(*d), 302)
}

type ToReferrer struct{}

func (d *ToReferrer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("Redirecting")
	defaultLogger.Debug(r.URL.Path)
	defaultLogger.Debug(r.Header.Get("Referer"))
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}
