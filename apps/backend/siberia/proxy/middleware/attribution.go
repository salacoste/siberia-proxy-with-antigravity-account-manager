package middleware

import (
	"net/http"

	"github.com/salacoste/siberia/siberia/proxy/privacy"
)

const (
	HeaderProvider = "x-antigravity-provider"
	HeaderModel    = "x-antigravity-model"
	HeaderAccount  = "x-antigravity-account"
)

// SetAttribution injects standard attribution headers into the response
func SetAttribution(w http.ResponseWriter, provider, model, account string) {
	w.Header().Set(HeaderProvider, provider)
	w.Header().Set(HeaderModel, model)
	w.Header().Set(HeaderModel, model)
	if account != "" {
		w.Header().Set(HeaderAccount, privacy.MaskEmail(account))
	}
}
