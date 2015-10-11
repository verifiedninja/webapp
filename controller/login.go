package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/verifiedninja/webapp/model"
	"github.com/verifiedninja/webapp/shared/passhash"
	"github.com/verifiedninja/webapp/shared/recaptcha"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"

	"github.com/gorilla/sessions"
)

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values["login_attempt"] == nil {
		sess.Values["login_attempt"] = 1
	} else {
		sess.Values["login_attempt"] = sess.Values["login_attempt"].(int) + 1
	}
}

// clearSessionVariables clears all the current session values
func clearSessionVariables(sess *sessions.Session) {
	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}
}

func LoginGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "login"
	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["login_attempt"] != nil && sess.Values["login_attempt"].(int) >= 5 {
		log.Println("Brute force login prevented")
		sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get database result
	result, err := model.UserByEmail(email)

	// Determine if user exists
	if err == sql.ErrNoRows {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values["login_attempt"]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.Status_id == 2 {
			// User inactive and display inactive message
			sess.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			sess.Save(r, w)
		} else if result.Status_id == 3 {
			// User email not verified
			sess.AddFlash(view.Flash{"You must confirm your email before login is allowed.", view.FlashWarning})
			sess.Save(r, w)
		} else if result.Status_id == 4 {
			// User email not re-verified
			sess.AddFlash(view.Flash{"You must re-confirm your email before login is allowed.", view.FlashWarning})
			sess.Save(r, w)
		} else if result.Status_id == 1 {
			// Log the login
			err = model.UserLoginCreate(int64(result.Id), r)
			if err != nil {
				log.Println(err)
			}

			// Login successfully
			clearSessionVariables(sess)
			sess.AddFlash(view.Flash{"Login successful!", view.FlashSuccess})
			sess.Values["id"] = result.Id
			sess.Values["email"] = email
			sess.Values["first_name"] = result.First_name
			sess.Values["last_name"] = result.Last_name
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			// User status unknown
			sess.AddFlash(view.Flash{"Account status is unknown so login is disabled.", view.FlashNotice})
			sess.Save(r, w)
		}
	} else {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values["login_attempt"]), view.FlashWarning})
		sess.Save(r, w)
	}

	// Show the login page again
	LoginGET(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		clearSessionVariables(sess)
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
