package middlewares

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type MWController struct {}

func (mw *MWController) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}

		if origin == "http://127.0.0.1:3000" {
			w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
		}

		//w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (mw *MWController) AccessLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		data := []string{r.Method, r.URL.String(), r.RemoteAddr, time.Now().UTC().String()}
		logrus.Info(strings.Join(data, " "))
	})
}
