package ujs

import "net/http"

func MethodRewriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if method := r.FormValue("_method"); method != "" {
			r.Method = method
		}
		next.ServeHTTP(w, r)
	})
}
