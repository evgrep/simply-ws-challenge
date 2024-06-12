package middleware

import (
	"net/http"
	"strings"
)

// AuthCheckHandler - represents authorization middleware. It chains a request requiring authorization.
type AuthCheckHandler struct {
	next IRequestHandler
}

// NewAuthCheckHandler - create a AuthCheckHandler. Accepts the next handler to call if authorization was
// successful.
func NewAuthCheckHandler(nextHandler IRequestHandler) *AuthCheckHandler {
	return &AuthCheckHandler{next: nextHandler}
}

func (h *AuthCheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// just pretend we issue tokens (stateful or stateless) we know how to validate
	const ValidDevToken = "8af4cc4fbf1eb641b14aeb7235bc7509"

	authHeader := strings.ToLower(r.Header.Get("Authorization"))
	authHeaderParts := strings.Split(authHeader, " ")

	if len(authHeaderParts) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bearerPart := strings.ToLower(authHeaderParts[0])

	if bearerPart != "bearer" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]
	if token != ValidDevToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if h.next != nil {
		h.next.Handle(w, r)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte{})
	}
}
