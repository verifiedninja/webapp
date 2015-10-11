package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/verifiedninja/webapp/model"
	emailer "github.com/verifiedninja/webapp/shared/email"
	"github.com/verifiedninja/webapp/shared/passhash"
	"github.com/verifiedninja/webapp/shared/random"
	"github.com/verifiedninja/webapp/shared/recaptcha"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"
)

func UserProfileGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the user photos
	photos, err := model.PhotosByUserId(uint64(sess.Values["id"].(uint32)))
	if err != nil {
		log.Println(err)
	}

	note := ""
	photo := ""
	status := uint8(0)
	date := time.Now()

	verified_private := false
	unverified_private := false
	rejected_private := false
	verified_public := false

	for _, v := range photos {
		if v.Initial == 1 {
			if v.Status_id == 1 {
				verified_private = true
			} else if v.Status_id == 2 {
				unverified_private = true
				note = v.Note
				photo = v.Path
				status = v.Status_id
				date = v.Updated_at
			} else if v.Status_id == 3 {
				rejected_private = true
				note = v.Note
				photo = v.Path
				status = v.Status_id
				date = v.Updated_at
			}
		} else {
			if v.Status_id == 1 {
				verified_public = true
			}
		}
	}

	user_id := strconv.Itoa(int(sess.Values["id"].(uint32)))

	// Display the view
	v := view.New(r)

	v.Vars["isNinja"] = false

	// If a private photo is verified, show the page
	if verified_private {
		v.Name = "user_profile"

		// Get the photo information
		imagesDB, err := model.PhotosByUserId(uint64(sess.Values["id"].(uint32)))
		if err != nil {
			log.Println(err)
			return
		}
		images := []Image{}
		for _, val := range imagesDB {
			img := Image{}
			img.Name = val.Path
			/*if val.Status_id == 1 {
				img.Path = "image/" + user_id + "/" + val.Path + ".jpg"
			} else {
				img.Path = photoPath + user_id + "/" + val.Path + ".jpg"
			}*/

			img.Path = "image/" + user_id + "/" + val.Path + ".jpg"

			img.Status_id = int(val.Status_id)
			img.Date = val.Updated_at.Format("Jan _2, 2006")
			img.Initial = int(val.Initial)
			img.Note = val.Note

			images = append(images, img)
		}
		v.Vars["images"] = images

		// Get the username information
		sites, err := model.UserinfoByUserId(uint64(sess.Values["id"].(uint32)))
		if err != nil {
			log.Println(err)
			return
		}
		for i, val := range sites {
			sites[i].Profile = strings.Replace(val.Profile, ":name", val.Username, -1)
		}
		v.Vars["sites"] = sites

		if len(sites) > 0 && verified_public {
			v.Vars["isNinja"] = true
		}

	} else {
		if unverified_private {
			// THIS NOTE MAY NOT BE FOR THE CORRECT PICTURE
			v.Vars["note"] = note
			//v.Vars["photo"] = photoPath + user_id + "/" + photo + ".jpg"
			v.Vars["photo"] = "image/" + user_id + "/" + photo + ".jpg"
			v.Vars["status_id"] = status
			v.Vars["date"] = date.Format("Jan _2, 2006")
			v.Vars["photo_id"] = photo
			v.Name = "user_unverified"
		} else if rejected_private {
			// THIS NOTE MAY NOT BE FOR THE CORRECT PICTURE
			v.Vars["note"] = note
			//v.Vars["photo"] = photoPath + user_id + "/" + photo + ".jpg"
			v.Vars["photo"] = "image/" + user_id + "/" + photo + ".jpg"
			v.Vars["status_id"] = status
			v.Vars["date"] = date.Format("Jan _2, 2006")
			v.Vars["photo_id"] = photo
			v.Name = "user_rejected"
		} else {
			http.Redirect(w, r, "/profile/initial", http.StatusFound)
			return
		}
	}

	v.Vars["first_name"] = sess.Values["first_name"]
	v.Render(w)
}

func UserSiteGET(w http.ResponseWriter, r *http.Request) {

	// Get session
	sess := session.Instance(r)

	// Does the user have a verified photo
	verified := isVerified(r)

	// Only allow access to this page if verified
	if verified {

		// Get database result
		sites, err := model.SiteList()
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}

		user_id := uint64(sess.Values["id"].(uint32))

		usernames, err := model.UsernamesByUserId(user_id)
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}

		// err == sql.ErrNoRows

		// Display the view
		v := view.New(r)
		v.Name = "user_site"

		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["sites"] = sites

		// Copy the usernames into a map so they can be used in the form inputs
		data := make(map[uint32]string)
		for _, u := range usernames {
			data[u.Site_id] = u.Name
		}
		v.Vars["data"] = data

		v.Render(w)
	} else {
		Error404(w, r)
	}
}

func UserSitePOST(w http.ResponseWriter, r *http.Request) {

	// Get session
	sess := session.Instance(r)

	// Does the user have a verified photo
	verified := isVerified(r)

	// Only allow access to this page if verified
	if verified {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}

		user_id := uint64(sess.Values["id"].(uint32))

		for site_id, values := range r.Form {

			site_idInt, err := strconv.Atoi(site_id)
			if err != nil {
				log.Println(err)
				continue
			}

			if len(values) < 1 {
				log.Println("Value is missing!", site_id)
				continue
			}

			username := values[0]

			if len(strings.TrimSpace(username)) < 1 {
				result, err := model.UsernameRemove(user_id, uint64(site_idInt))

				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"There was an issue with the username: " + username + ". Please try again later.", view.FlashError})
					sess.Save(r, w)
				} else {
					affected, err := result.RowsAffected()
					if err != nil {
						log.Println(err)
					} else if affected > 0 {
						sess.AddFlash(view.Flash{"Removed username", view.FlashSuccess})
						sess.Save(r, w)
					}
				}
			} else {
				err = model.UsernameAdd(user_id, username, uint64(site_idInt))
				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"There was an issue with the username: " + username + ". Please try again later.", view.FlashError})
					sess.Save(r, w)
				} else {
					sess.AddFlash(view.Flash{"Saved username: " + username, view.FlashSuccess})
					sess.Save(r, w)
				}
			}
		}

		//UserSiteGET(w, r)
		http.Redirect(w, r, "/profile", http.StatusFound)
	} else {
		Error404(w, r)
	}
}

func UserInformationGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	user_id := uint64(sess.Values["id"].(uint32))

	demo, err := model.DemographicByUserId(user_id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}

	ethnicity, err := model.EthnicityByUserId(user_id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}

	e := make(map[string]int)

	for i := 0; i <= 9; i++ {
		ee := fmt.Sprintf("%v", i)
		e["E"+ee] = 0
		for j := 0; j < len(ethnicity); j++ {
			if int(ethnicity[j].Type_id) == i {
				e["E"+ee] = 1
			}
		}
	}

	// Display the view
	v := view.New(r)
	v.Name = "user_info"
	//v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]

	v.Vars["demographic"] = demo
	v.Vars["ethnicity"] = e

	v.Render(w)
}

func UserInformationPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"birth_month", "birth_day", "birth_year", "gender", "height_feet", "height_inches", "ethnicity"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}

	user_id := uint64(sess.Values["id"].(uint32))

	d := model.Demographic{}

	// Get form values
	bm, err := strconv.Atoi(r.FormValue("birth_month"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}
	bd, _ := strconv.Atoi(r.FormValue("birth_day"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}
	by, _ := strconv.Atoi(r.FormValue("birth_year"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}

	d.Birth_month = uint8(bm)
	d.Birth_day = uint8(bd)
	d.Birth_year = uint16(by)

	d.Gender = r.FormValue("gender")

	hf, _ := strconv.Atoi(r.FormValue("height_feet"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}

	d.Height_feet = uint8(hf)

	hi, _ := strconv.Atoi(r.FormValue("height_inches"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}

	d.Height_inches = uint8(hi)

	we, _ := strconv.Atoi(r.FormValue("weight"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}
	d.Weight = uint16(we)

	err = model.DemographicAdd(user_id, d)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}

	err = model.EthnicityAdd(user_id, r.Form["ethnicity"])
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserInformationGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"Settings saved.", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func UserEmailGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	user_id := int64(sess.Values["id"].(uint32))
	if !isVerifiedEmail(r, user_id) {
		sess.AddFlash(view.Flash{"You can't change you email again until you verify your current email.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// Display the view
	v := view.New(r)
	v.Name = "user_email"
	v.Vars["emailold"] = sess.Values["email"]

	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)

	v.Render(w)
}

func UserEmailPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	user_id := int64(sess.Values["id"].(uint32))
	if !isVerifiedEmail(r, user_id) {
		sess.AddFlash(view.Flash{"You can't change you email again until you verify your current email.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"email"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		UserEmailGET(w, r)
		return
	}

	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		UserEmailGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	emailOld := sess.Values["email"]

	if email == emailOld {
		sess.AddFlash(view.Flash{"New email cannot be the same as the old email.", view.FlashError})
		sess.Save(r, w)
		UserEmailGET(w, r)
		return
	}

	// Get database result
	err := model.UserEmailUpdate(user_id, email)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			sess.AddFlash(view.Flash{"That email already exists in the database. Please use a different one.", view.FlashError})
		} else {
			// Display error message
			log.Println(err)
			sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		}

		sess.Save(r, w)
		UserEmailGET(w, r)
		return
	}

	first_name := fmt.Sprintf("%v", sess.Values["first_name"])

	// Create the email verification string
	md := random.Generate(32)

	// Add the hash to the database
	err = model.UserEmailVerificationCreate(user_id, md)
	if err != nil {
		log.Println(err)
	}
	err = model.UserReverify(user_id)
	if err != nil {
		log.Println(err)
	}

	c := view.ReadConfig()

	// Email the hash to the user
	err = emailer.SendEmail(email, "Email Verification for Verified.ninja", "Hi "+first_name+",\n\nTo verify your email address ("+email+"), please click here: "+c.BaseURI+"emailverification/"+md)
	if err != nil {
		log.Println(err)
	}

	// Login successfully
	sess.AddFlash(view.Flash{"Email updated! You must verify your email before you can login again.", view.FlashSuccess})
	sess.Values["email"] = email
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func UserPasswordGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "user_password"

	v.Render(w)
}

func UserPasswordPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"passwordold", "passwordnew"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		UserPasswordGET(w, r)
		return
	}

	user_id := int64(sess.Values["id"].(uint32))

	// Form values
	passwordOld := r.FormValue("passwordold")

	passwordNew, errp := passhash.HashString(r.FormValue("passwordnew"))

	if passwordOld == r.FormValue("passwordnew") {
		sess.AddFlash(view.Flash{"New password cannot be the same as the old password.", view.FlashError})
		sess.Save(r, w)
		UserPasswordGET(w, r)
		return
	}

	// If password hashing failed
	if errp != nil {
		log.Println(errp)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserPasswordGET(w, r)
		return
	}

	// Get database result
	result, err := model.UserByEmail(fmt.Sprintf("%v", sess.Values["email"]))

	// Determine if user exists
	if err == sql.ErrNoRows {
		sess.AddFlash(view.Flash{"Password is incorrect.", view.FlashWarning})
		sess.Save(r, w)
		UserPasswordGET(w, r)
		return
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(r, w)
		UserPasswordGET(w, r)
		return
	} else if passhash.MatchString(result.Password, passwordOld) {

		err = model.UserPasswordUpdate(user_id, passwordNew)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
			sess.Save(r, w)
			UserPasswordGET(w, r)
			return
		}

		// Password matches
		sess.AddFlash(view.Flash{"Password changed!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		sess.AddFlash(view.Flash{"Password is incorrect.", view.FlashWarning})
		sess.Save(r, w)
		UserPasswordGET(w, r)
		return
	}
}

/*******************************************************************************
Helpers
*******************************************************************************/

func isVerified(r *http.Request) bool {
	// Get session
	sess := session.Instance(r)

	// Get the user photos
	photos, err := model.PhotosByUserId(uint64(sess.Values["id"].(uint32)))
	if err != nil {
		log.Println(err)
	}

	verified := false

	for _, v := range photos {
		if v.Status_id == 1 {
			verified = true
			break
		}
	}

	return verified
}

func isVerifiedPrivate(r *http.Request, user_id uint64) bool {
	// Get the user photos
	photos, err := model.PhotosByUserId(user_id)
	if err != nil {
		log.Println(err)
	}

	verified := false

	for _, v := range photos {
		if v.Status_id == 1 && v.Initial == 1 {
			verified = true
			break
		}
	}

	return verified
}

func isVerifiedPublic(r *http.Request, user_id uint64) bool {
	// Get the user photos
	photos, err := model.PhotosByUserId(user_id)
	if err != nil {
		log.Println(err)
	}

	verified := false

	for _, v := range photos {
		if v.Status_id == 1 && v.Initial == 0 {
			verified = true
			break
		}
	}

	return verified
}

func isVerifiedEmail(r *http.Request, user_id int64) bool {
	// Get the user photos
	user, err := model.UserStatusByUserId(user_id)
	if err != nil {
		log.Println(err)
	}

	if user.Status_id == 1 {
		return true
	}

	return false
}
