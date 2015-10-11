package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/verifiedninja/webapp/model"
	emailer "github.com/verifiedninja/webapp/shared/email"
	"github.com/verifiedninja/webapp/shared/random"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func PublicUsernameGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	site := params.ByName("site")
	username := params.ByName("username")

	// Display the view
	v := view.New(r)

	v.Vars["isSelf"] = false
	v.Vars["verified_email"] = false

	user_info, err := model.UserByUsername(username, site)
	if err == sql.ErrNoRows {
		v.Vars["verified_private"] = false
		v.Vars["verified_public"] = false
		v.Vars["exists"] = false
	} else if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	} else {

		v.Vars["verified_email"] = isVerifiedEmail(r, int64(user_info.Id))

		v.Vars["exists"] = true

		if sess.Values["id"] != nil {
			if sess.Values["id"] == user_info.Id {
				v.Vars["isSelf"] = true
			}
		}

		if isVerifiedPublic(r, uint64(user_info.Id)) && isVerifiedPrivate(r, uint64(user_info.Id)) {
			v.Vars["verified_public"] = true

			// Get the photo information
			//user_id := strconv.Itoa(int(sess.Values["id"].(uint32)))
			user_id_string := strconv.Itoa(int(user_info.Id))
			imagesDB, err := model.PhotosByUserId(uint64(user_info.Id))

			if err != nil {
				log.Println(err)
				return
			}
			images := []Image{}
			for _, val := range imagesDB {
				img := Image{}
				img.Name = val.Path
				/*if val.Status_id == 1 {
					img.Path = "image/" + user_id_string + "/" + val.Path + ".jpg"
				} else {
					img.Path = photoPath + user_id_string + "/" + val.Path + ".jpg"
				}*/

				img.Path = "image/" + user_id_string + "/" + val.Path + ".jpg"

				img.Status_id = int(val.Status_id)
				img.Date = val.Updated_at.Format("Jan _2, 2006")

				// Only allows verified images right now
				if val.Status_id == 1 && val.Initial == 0 {
					images = append(images, img)
				}
			}
			v.Vars["site"] = user_info.Site
			v.Vars["profile"] = strings.Replace(user_info.Profile, ":name", user_info.Username, -1)

			v.Vars["images"] = images

		} else if isVerifiedPrivate(r, uint64(user_info.Id)) {
			v.Vars["verified_private"] = true
		} else {
			v.Vars["verified_private"] = false
		}
	}

	v.Name = "public_username"
	v.Vars["username"] = username
	//v.Vars["site"] = user_info.Site
	//v.Vars["profile"] = user_info.Profile
	v.Vars["home"] = user_info.Home
	v.Render(w)
}

func APIVerifyUserGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	//sess := session.Instance(r)

	user_id := uint64(0)
	other_user_id := uint64(0)

	userkey := r.URL.Query().Get("userkey")
	token := r.URL.Query().Get("token")

	auth_info, err := model.ApiAuthenticationByKeys(userkey, token)
	if err == sql.ErrNoRows {
		Error401(w, r)
		return
	} else if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	// If the user is logged in
	/*if sess.Values["id"] != nil {
		user_id = uint64(sess.Values["id"].(uint32))
	}*/

	user_id = uint64(auth_info.User_id)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	site := params.ByName("site")
	username := params.ByName("username")

	vn := VerifiedNinja{}

	user_info, err := model.UserByUsername(username, site)
	if err == sql.ErrNoRows {
	} else if err != nil {
		log.Println(err)
	} else {

		other_user_id = uint64(user_info.Id)

		// Get the user photos
		photos, err := model.PhotosByUserId(uint64(user_info.Id))
		if err != nil {
			log.Println(err)
		}

		for _, v := range photos {
			if v.Initial == 1 {
				if v.Status_id == 1 {
					vn.PrivatePhotoVerified = true
				}
			} else {
				if v.Status_id == 1 {
					vn.PublicPhotoVerified = true
				}
			}
		}

		// If a private photo is verified, show the page
		if vn.PrivatePhotoVerified && vn.PublicPhotoVerified {

			// Get the username information
			sites, err := model.UserinfoByUserId(uint64(user_info.Id))
			if err != nil {
				log.Println(err)
			} else {
				for _, s := range sites {
					if strings.ToLower(s.Site) == strings.ToLower(site) {
						vn.RegisteredUsername = true
						vn.VerifiedNinja = true
						break
					}
				}
			}
		}
	}

	//log.Println("API Check - is Ninja?:", username, site, vn.VerifiedNinja)

	err = model.TrackRequestAPI(user_id, r, other_user_id, vn.VerifiedNinja)
	if err != nil {
		log.Println(err)
	}

	js, err := json.Marshal(vn)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type VerifiedNinja struct {
	RegisteredUsername   bool `json:RegisteredUsernames`
	PrivatePhotoVerified bool `json:PrivatePhotoVerified`
	PublicPhotoVerified  bool `json:PublicPhotoVerified`
	VerifiedNinja        bool `json:VerifiedNinja`
}

func CronNotifyEmailExpireGET(w http.ResponseWriter, r *http.Request) {
	result, err := model.EmailsWithVerificationIn30Days()

	if err != nil {
		log.Println(err)
		Error500(w, r)
	}

	c := view.ReadConfig()

	for _, v := range result {
		//log.Println(v.First_name, " expires on ", v.Expiring, v.Expired, v.Updated_at)

		user_id := int64(v.Id)
		email := v.Email
		first_name := v.First_name

		if v.Expiring {
			// Create the email verification string
			md := random.Generate(32)

			// Add the hash to the database
			err = model.UserEmailVerificationCreate(user_id, md)
			if err != nil {
				log.Println(err)
			}

			// Email the hash to the user
			err = emailer.SendEmail(email, "Email Verification Required for Verified.ninja", "Hi "+first_name+",\n\nTo keep your account active, please verify your email address by clicking on this link: "+c.BaseURI+"emailverification/"+md+"\n\nYour account will expire in 5 days if you don't verify your email.")
			if err != nil {
				log.Println(err)
			}
		} else if v.Expired {
			err = model.UserReverify(user_id)
			if err != nil {
				log.Println(err)
			}

			user_info, err := model.EmailVerificationTokenByUserId(uint64(user_id))
			if err != nil {
				log.Println(err)
			}

			md := user_info.Token

			// Email the hash to the user
			err = emailer.SendEmail(email, "Account Locked on Verified.ninja", "Hi "+first_name+",\n\nIt's been over 30 days since you verified your email. To unlock your account, please verify your email address by clicking on this link: "+c.BaseURI+"emailverification/"+md)
			if err != nil {
				log.Println(err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"Done": true}`))
}

type ChromeInfo struct {
	Userkey string
	Token   string
}

func APIRegisterChromeGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If the user is logged in
	if sess.Values["id"] == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Authentication required"}`))
		return
	}

	user_id := uint64(sess.Values["id"].(uint32))

	ci := ChromeInfo{}

	auth, err := model.ApiAuthenticationByUserId(user_id)
	// If there is no record yet, create one
	if err == sql.ErrNoRows {

		// Create values
		ci.Userkey = random.Generate(32)
		ci.Token = random.Generate(32)

		err := model.ApiAuthenticationCreate(user_id, ci.Userkey, ci.Token)
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}
	} else if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	} else {
		ci.Userkey = auth.Userkey
		ci.Token = auth.Token
	}

	js, err := json.Marshal(ci)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
