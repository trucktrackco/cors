package cors

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const options = "OPTIONS"

func addCORSHeader(env interface{}, w *http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		if strings.Index(origin, "http://") == -1 && (env == "dev") {
			log.Println("Dev environment detected", origin)
			origin = "http://" + origin
		}

		if strings.Index(origin, "https://") == -1 && (env == "stage" || env == "prod") {
			log.Println("Stage/prod environment detected", origin)
			origin = "https://" + origin
		}
	}
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Header.Get("Access-Control-Request-Method") != "" && r.Method == options {
		(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		(*w).Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Origin,Content-Type, Accept, Authorization")
	}
}

// OriginFunc ...
func OriginFunc(env interface{}, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addCORSHeader(env, &w, r)
		if r.Method != options {
			h.ServeHTTP(w, r)
			fmt.Printf("%v %s %s \n", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL.Path)
		}
	})
}

// Origin ...
func Origin(env interface{}, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addCORSHeader(env, &w, r)
		if r.Method != options {
			h.ServeHTTP(w, r)
			fmt.Printf("%v %s %s \n", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL.Path)
		}
	})
}
