package middlewares

import (
	"net/http"
)

type MWController struct {}

func (mw *MWController) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		//w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		next.ServeHTTP(w, r)
	})
}
