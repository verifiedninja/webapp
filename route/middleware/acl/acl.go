package acl

import (
	"log"
	"net/http"

	"github.com/verifiedninja/webapp/controller"
	"github.com/verifiedninja/webapp/model"
	"github.com/verifiedninja/webapp/shared/session"
)

// DisallowAuth does not allow authenticated users to access the page
func DisallowAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is authenticated, don't allow them to access the page
		if sess.Values["id"] != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// DisallowAnon does not allow anonymous users to access the page
func DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// DisallowAnon does not allow anonymous users to access the page
func AllowOnlyAdministrator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current user role
		currentUID, isLoggedIn := controller.CurrentUserId(r)

		// If user is not authenticated, don't allow them to access the page
		if !isLoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Get the role
		role, err := model.RoleByUserId(int64(currentUID))
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Only allow Administrators
		if role.Level_id != model.Role_level_Administrator {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}
