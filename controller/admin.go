package controller

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/verifiedninja/webapp/model"
	emailer "github.com/verifiedninja/webapp/shared/email"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	Id              int
	FirstName       string
	LastName        string
	Images          []Image
	Token           string
	VerifiedCount   int
	UnverifiedCount int
	SiteCount       int
	Ninja           bool
	Email           bool
}

type Image struct {
	Id        int
	Name      string
	Path      string
	Date      string
	Status_id int
	Initial   int
	Note      string
}

// Displays the default home page
func AdminGET(w http.ResponseWriter, r *http.Request) {
	dirs, err := filepath.Glob(photoPath + "*")
	if err != nil {
		log.Println(err)
	}

	users := []User{}

	ds := string(os.PathSeparator)

	for _, v := range dirs {
		u := User{}
		idRaw := v[strings.LastIndex(v, ds)+1:]
		u.Id, err = strconv.Atoi(idRaw)
		if err != nil {
			log.Println(err)
			continue
		}

		info, err := model.UserNameById(u.Id)
		if err == sql.ErrNoRows {
			log.Println("User is not found in database:", u.Id)
			continue
		} else if err != nil {
			log.Println(err)
			continue
		}

		u.FirstName = info.First_name
		u.LastName = info.Last_name

		privateVerifiedCount := 0
		publicVerifiedCount := 0

		// Get the photo information
		user_id := strconv.Itoa(u.Id)
		imagesDB, err := model.PhotosByUserId(uint64(u.Id))
		if err != nil {
			log.Println(err)
			return
		}
		//images := []Image{}
		for _, val := range imagesDB {
			img := Image{}
			img.Name = val.Path
			if val.Status_id == 1 {
				u.VerifiedCount += 1

				if val.Initial == 1 {
					privateVerifiedCount += 1
				} else if val.Initial == 0 {
					publicVerifiedCount += 1
				}

			} else if val.Status_id == 2 {
				u.UnverifiedCount += 1
			}

			img.Path = "image/" + user_id + "/" + val.Path + ".jpg"

			img.Status_id = int(val.Status_id)
			img.Date = val.Updated_at.Format("Jan _2, 2006")
			img.Initial = int(val.Initial)
			u.Images = append(u.Images, img)
		}

		// Get the user verification code
		token_info, err := model.UserTokenByUserId(uint64(u.Id))
		if err == sql.ErrNoRows {
			log.Println(err)
			token_info.Token = "TOKEN IS MISSING"
		} else if err != nil {
			log.Println(err)
			token_info.Token = "TOKEN IS MISSING"
		}
		u.Token = token_info.Token

		// Get the username information
		sites, err := model.UserinfoByUserId(uint64(u.Id))
		if err != nil {
			log.Println(err)
			return
		}
		u.SiteCount = len(sites)

		u.Email = isVerifiedEmail(r, int64(u.Id))

		if u.SiteCount > 0 && privateVerifiedCount > 0 && publicVerifiedCount > 0 && u.Email {
			u.Ninja = true
		}

		users = append(users, u)
	}

	// Display the view
	v := view.New(r)
	v.Name = "admin"
	v.Vars["users"] = users
	v.Render(w)
}

// Displays the default home page
func AdminUserGET(w http.ResponseWriter, r *http.Request) {
	var params = context.Get(r, "params").(httprouter.Params)
	userid := params.ByName("userid")
	user_id, _ := strconv.Atoi(userid)

	users := []User{}

	for _, v := range []int{user_id} {
		u := User{}
		u.Id = v

		info, err := model.UserNameById(u.Id)
		if err == sql.ErrNoRows {
			log.Println("User is not found in database:", u.Id)
			continue
		} else if err != nil {
			log.Println(err)
			continue
		}

		u.FirstName = info.First_name
		u.LastName = info.Last_name

		// Get the photo information
		user_id := strconv.Itoa(u.Id)
		imagesDB, err := model.PhotosByUserId(uint64(u.Id))
		if err != nil {
			log.Println(err)
			return
		}
		//images := []Image{}
		for _, val := range imagesDB {
			img := Image{}
			img.Name = val.Path
			img.Path = "image/" + user_id + "/" + val.Path + ".jpg"

			img.Status_id = int(val.Status_id)
			img.Date = val.Updated_at.Format("Jan _2, 2006")
			img.Initial = int(val.Initial)
			u.Images = append(u.Images, img)
		}

		// Get the user verification code
		token_info, err := model.UserTokenByUserId(uint64(u.Id))
		if err == sql.ErrNoRows {
			log.Println(err)
			token_info.Token = "TOKEN IS MISSING"
		} else if err != nil {
			log.Println(err)
			token_info.Token = "TOKEN IS MISSING"
		}
		u.Token = token_info.Token
		users = append(users, u)
	}

	// Display the view
	v := view.New(r)
	v.Name = "admin_all"
	v.Vars["users"] = users
	v.Render(w)
}

// Displays the default home page
func AdminAllGET(w http.ResponseWriter, r *http.Request) {
	dirs, err := filepath.Glob(photoPath + "*")
	if err != nil {
		log.Println(err)
	}

	users := []User{}

	ds := string(os.PathSeparator)

	for _, v := range dirs {
		u := User{}
		idRaw := v[strings.LastIndex(v, ds)+1:]
		u.Id, err = strconv.Atoi(idRaw)
		if err != nil {
			log.Println(err)
			continue
		}

		info, err := model.UserNameById(u.Id)
		if err == sql.ErrNoRows {
			log.Println("User is not found in database:", u.Id)
			continue
		} else if err != nil {
			log.Println(err)
			continue
		}

		u.FirstName = info.First_name
		u.LastName = info.Last_name

		/*files, err := filepath.Glob(photoPath + idRaw + "/*")
		if err != nil {
			log.Println(err)
			continue
		}

		for _, v := range files {
			i := Image{}
			i.Name = v[strings.LastIndex(v, ds)+1:]
			iid, _ := strconv.Atoi(strings.Replace(i.Name, `.jpg`, ``, -1))
			i.Id = iid
			i.Path = strings.Replace(v, `\`, `/`, -1)
			u.Images = append(u.Images, i)
		}*/

		// Get the photo information
		user_id := strconv.Itoa(u.Id)
		imagesDB, err := model.PhotosByUserId(uint64(u.Id))
		if err != nil {
			log.Println(err)
			return
		}
		//images := []Image{}
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
			u.Images = append(u.Images, img)
		}

		//uid := sess.Values["id"].(uint32)

		// Get the user verification code
		token_info, err := model.UserTokenByUserId(uint64(u.Id))
		if err == sql.ErrNoRows {
			log.Println(err)
			token_info.Token = "TOKEN IS MISSING"
		} else if err != nil {
			log.Println(err)
			token_info.Token = "TOKEN IS MISSING"
		}
		u.Token = token_info.Token
		users = append(users, u)
	}

	// Display the view
	v := view.New(r)
	v.Name = "admin_all"
	v.Vars["users"] = users
	v.Render(w)
}

func AdminApproveGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	userid := params.ByName("userid")
	picid := params.ByName("picid")
	uid, _ := strconv.Atoi(userid)

	err := model.PhotoApprove(picid, uint64(uid))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Photo approved!", view.FlashSuccess})
		sess.Save(r, w)

		user_info, err := model.UserEmailByUserId(int64(uid))
		if err != nil {
			log.Println()
		} else {
			c := view.ReadConfig()

			// Email the update to the user
			err := emailer.SendEmail(user_info.Email, "Photo Approved on Verified.ninja",
				"Hi "+user_info.First_name+",\n\nYour photo ("+picid+") was approved!\n\nLogin to see your updated profile: "+c.BaseURI)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Display the view
	v := view.New(r)
	v.SendFlashes(w)
}

func AdminRejectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)

	userid := params.ByName("userid")
	picid := params.ByName("picid")
	note := r.FormValue("note")
	uid, _ := strconv.Atoi(userid)

	err := model.PhotoReject(picid, uint64(uid), note)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Photo rejected!", view.FlashSuccess})
		sess.Save(r, w)

		user_info, err := model.UserEmailByUserId(int64(uid))
		if err != nil {
			log.Println()
		} else {
			c := view.ReadConfig()

			// Email the update to the user
			err := emailer.SendEmail(user_info.Email, "Photo Rejected on Verified.ninja",
				"Hi "+user_info.First_name+",\n\nYour photo ("+picid+") was rejected for the following reason(s):\n"+note+"\n\nPlease upload a new private photo for verification: "+c.BaseURI)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Display the view
	v := view.New(r)
	v.SendFlashes(w)
}

func AdminUnverifyGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	userid := params.ByName("userid")
	picid := params.ByName("picid")
	uid, _ := strconv.Atoi(userid)

	err := model.PhotoUnverify(picid, uint64(uid))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Photo unverified!", view.FlashSuccess})
		sess.Save(r, w)
	}

	// Display the view
	v := view.New(r)
	v.SendFlashes(w)
}
