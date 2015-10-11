package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/verifiedninja/webapp/model"
	emailer "github.com/verifiedninja/webapp/shared/email"
	"github.com/verifiedninja/webapp/shared/recaptcha"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"
)

// Displays the default home page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about"
	v.Render(w)
}

// Displays the default home page
func TermsGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "terms"
	v.Render(w)
}

func ContactGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "contact"

	// If the user is logged in
	if sess.Values["id"] != nil {
		// Refill any form fields
		view.Repopulate([]string{"message"}, r.Form, v.Vars)
		v.Vars["email"] = sess.Values["email"]
		v.Vars["fullname"] = fmt.Sprintf("%v %v", sess.Values["first_name"], sess.Values["last_name"])
		v.Vars["logged_in"] = true
	} else {
		// Refill any form fields
		view.Repopulate([]string{"email", "fullname", "message"}, r.Form, v.Vars)
	}

	v.Render(w)
}

func ContactPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"email", "fullname", "message"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		ContactGET(w, r)
		return
	}

	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		ContactGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	name := r.FormValue("fullname")
	message := r.FormValue("message")
	ip, err := model.GetRemoteIP(r)
	if err != nil {
		log.Println(err)
	}

	user := "Guest"

	if sess.Values["id"] != nil {
		user = fmt.Sprintf("Registered (%v)", sess.Values["id"])
	}

	// Email the hash to the user
	err = emailer.SendEmail(emailer.ReadConfig().From, "Contact Submission for Verified.ninja", "From: "+
		name+" <"+email+">\nUser: "+user+"\nIP: "+ip+"\nMessage: "+message)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		ContactGET(w, r)
		return
	}

	// Post successful
	sess.AddFlash(view.Flash{"Thanks for the message! We'll get back to you in a bit.", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
