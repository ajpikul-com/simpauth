package simpauth

import (
	"net/http"

	// "github.com/pkg/errors"

	"github.com/ajpikul-com/ilog"
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

type Bouncer struct {
	base http.Handler
}

func (h *Bouncer) Init(handler http.Handler) error {
	h.base = handler
	return nil
}

func (h *Bouncer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.base.ServeHTTP(w, r)
}
