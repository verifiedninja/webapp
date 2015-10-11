package controller

import (
	"net/http"

	"github.com/verifiedninja/webapp/shared/view"
)

// Displays the default home page
func VerifyUsernameGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "verify_username"
	v.Render(w)
}
