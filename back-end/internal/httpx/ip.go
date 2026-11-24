package httpx

import (
	"net"
	"net/http"
	"strings"
)

// ExtratorIP descobre o IP do cliente. Quando ConfiarProxy é verdadeiro (ex.:
// atrás do proxy do Render) usa o primeiro valor de X-Forwarded-For.
type ExtratorIP struct {
	ConfiarProxy bool
}

func (e ExtratorIP) De(r *http.Request) string {
	if e.ConfiarProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return strings.TrimSpace(strings.Split(xff, ",")[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
