package logrequest

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/verifiedninja/webapp/model"
	"github.com/verifiedninja/webapp/shared/session"
)

// Handler will log the HTTP requests
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// Handler will log the HTTP requests
func Database(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.URL.Path, "/api/v1/verify") {
			// Get session
			sess := session.Instance(r)

			user_id := uint64(0)

			// If the user is logged in
			if sess.Values["id"] != nil {
				user_id = uint64(sess.Values["id"].(uint32))
			}

			err := model.TrackRequestURL(user_id, r)
			if err != nil {
				log.Println(err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
