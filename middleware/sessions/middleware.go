package sessions

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func Use(next http.Handler) http.Handler {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nr := r.WithContext(context.WithValue(r.Context(), sessionStoreKey, store))
		next.ServeHTTP(rw, nr)
	})
}

func Store(r *http.Request) sessions.Store {
	store, _ := r.Context().Value(sessionStoreKey).(sessions.Store)
	return store
}

type sessionStoreKeyType struct{}

var sessionStoreKey = sessionStoreKeyType{}
