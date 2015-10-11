package controller

import (
	"fmt"
	"log"
	//"log"
	"net/http"
	"strings"

	"github.com/verifiedninja/webapp/model"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"
)

// Displays the default home page
func Index(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If the user is logged in
	if sess.Values["id"] != nil {

		// Get the current user role
		currentUID, _ := CurrentUserId(r)
		role, err := model.RoleByUserId(int64(currentUID))
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}

		if role.Level_id == model.Role_level_User {
			http.Redirect(w, r, "/profile", http.StatusFound)
			return
		} else {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}

	} else {
		// Display the view
		v := view.New(r)
		v.Name = "anon_home"
		v.Render(w)
	}
}

// Error401 handles 401 - Unauthorized
func Error401(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	//fmt.Fprint(w, "Unauthorized 401")

	// Display the view
	v := view.New(r)
	v.Name = "error_401"
	v.Render(w)
}

// Error404 handles 404 - Page Not Found
func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	//fmt.Fprint(w, "Not Found 404")

	// Display the view
	v := view.New(r)
	v.Name = "error_404"
	v.Render(w)
}

// Error500 handles 500 - Internal Server Error
func Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	//fmt.Fprint(w, "Internal Server Error 500")

	// Display the view
	v := view.New(r)
	v.Name = "error_500"
	v.Render(w)
}

// InvalidToken handles CSRF attacks
func InvalidToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, `Your token <strong>expired</strong>, click <a href="javascript:void(0)" onclick="window.history.back()">here</a> to try again.`)
}

// Static maps static files
func Static(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	if strings.HasSuffix(r.URL.Path, "/") {
		Error404(w, r)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}
