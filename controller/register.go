package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/verifiedninja/webapp/model"
	emailer "github.com/verifiedninja/webapp/shared/email"
	"github.com/verifiedninja/webapp/shared/passhash"
	"github.com/verifiedninja/webapp/shared/pushover"
	"github.com/verifiedninja/webapp/shared/random"
	"github.com/verifiedninja/webapp/shared/recaptcha"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

var (
	captchaEnabled = true
)

type CaptchaResponse struct {
	success    bool   `json:"success"`
	errorCodes string `json:"error-codes"`
}

func RegisterGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "register"
	// Refill any form fields
	view.Repopulate([]string{"first_name", "last_name", "email"}, r.Form, v.Vars)
	v.Render(w)
}

func EmailVerificationGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	token := params.ByName("token")

	result, err := model.UserEmailVerified(token)

	// Will only error if there is a problem with the query
	if err != nil {
		//log.Println(err)
		sess.AddFlash(view.Flash{"Email verification link is not valid. Please check the latest email we sent you.", view.FlashError})
		sess.Save(r, w)
	} else if result {
		sess.AddFlash(view.Flash{"Email verified. You may login now.", view.FlashSuccess})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		log.Println("Brute force register prevented")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"first_name", "last_name", "email", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}

	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}

	// Get form values
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	email := r.FormValue("email")

	password, errp := passhash.HashString(r.FormValue("password"))

	// If password hashing failed
	if errp != nil {
		log.Println(errp)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	// Get database result
	_, err := model.UserIdByEmail(email)

	if err == sql.ErrNoRows { // If success (no user exists with that email)
		result, ex := model.UserCreate(first_name, last_name, email, password)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
		} else {

			// Create the email verification string
			md := random.Generate(32)

			// Get the user ID
			user_id, _ := result.LastInsertId()

			// Add the user role
			model.RoleCreate(user_id, model.Role_level_User)

			// Add the hash to the database
			model.UserEmailVerificationCreate(user_id, md)

			c := view.ReadConfig()

			// Email the hash to the user
			err := emailer.SendEmail(email, "Email Verification for Verified.ninja", "Hi "+first_name+",\n\nTo verify your email address, please click here: "+c.BaseURI+"emailverification/"+md)
			if err != nil {
				log.Println(err)
			}

			// TODO This is just temporary for testing
			log.Println("Email Verification Link:", c.BaseURI+"emailverification/"+md)

			po, err := pushover.New()
			if err == pushover.ErrPushoverDisabled {
				// Nothing
			} else if err != nil {
				log.Println(err)
			} else {
				err = po.Message(first_name + " " + last_name + "(" + fmt.Sprintf("%v", user_id) + ") created an account. You can view the account here:\nhttps://verified.ninja/admin/user/" + fmt.Sprintf("%v", user_id))
				if err != nil {
					log.Println(err)
				}
			}

			sess.AddFlash(view.Flash{"Account created successfully for: " + email + ". Please click the verification link in your email.", view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	} else if err != nil { // Catch all other errors
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else { // Else the user already exists
		sess.AddFlash(view.Flash{"Account already exists for: " + email, view.FlashError})
		sess.Save(r, w)
	}

	// Display the page
	RegisterGET(w, r)
}
