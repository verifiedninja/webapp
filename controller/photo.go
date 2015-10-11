package controller

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/verifiedninja/webapp/model"
	fs "github.com/verifiedninja/webapp/shared/filesystem"
	"github.com/verifiedninja/webapp/shared/photo"
	"github.com/verifiedninja/webapp/shared/pushover"
	"github.com/verifiedninja/webapp/shared/random"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

var (
	supportedFileTypes = []string{"image/jpeg", "image/jpg", "image/gif", "image/png"}
	photoLimit         = 10
	photoPath          = "private/images/"
)

type Ethnicity_type struct {
	E1 uint8
	E2 uint8
	E3 uint8
	E4 uint8
	E5 uint8
	E6 uint8
	E7 uint8
	E8 uint8
	E9 uint8
}

func renderImage(w http.ResponseWriter, r *http.Request, userid, picid string, mark bool) (*bytes.Buffer, error) {
	filename := photoPath + userid + "/" + picid

	file_info, err := photo.FileDimensions(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	water := "static/resources/verifiedninja32.png"
	padding := 5

	if file_info.Height >= 800 || file_info.Width >= 800 {
		water = "static/resources/verifiedninja.png"
		padding = 20
	}

	watermarkFile, err := os.Open(water)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer watermarkFile.Close()
	watermarkImage, _, err := image.Decode(watermarkFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Bottom Right
	//offset := image.Pt(img.Bounds().Size().X-watermarkImage.Bounds().Size().X-padding, img.Bounds().Size().Y-32-padding)

	// Top Right
	offset := image.Pt(img.Bounds().Size().X-watermarkImage.Bounds().Size().X-padding, padding)

	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)

	if mark {
		draw.Draw(m, watermarkImage.Bounds().Add(offset), watermarkImage, image.ZP, draw.Over)
	}

	// Convert image to JPEG
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, m, nil); err != nil {
		log.Println("unable to encode image.")
		return nil, err
	}

	return buffer, nil
}

// Returns the user_id in uint32 and whether the a user is logged in
func CurrentUserId(r *http.Request) (uint32, bool) {
	// Get session
	sess := session.Instance(r)

	if sess.Values["id"] == nil {
		return 0, false
	}

	return sess.Values["id"].(uint32), true
}

// photoAccessAllowed returns: 1 - Is use allows access 2 - Should photo be marked
func photoAccessAllowed(r *http.Request, user_id uint64, pic_id string) (bool, bool, error) {
	// Get the photo info

	photoInfo, err := model.PhotoInfoByPath(user_id, strings.Replace(pic_id, ".jpg", "", -1))
	if err != nil {
		return false, false, err
	}

	// Get the current user role
	role_level := uint8(0)
	currentUID, loggedIn := CurrentUserId(r)
	if loggedIn {
		role, err := model.RoleByUserId(int64(currentUID))
		if err != nil {
			return false, false, err
		}
		role_level = role.Level_id
	}

	// Check if the current user has access to the photo
	if (photoInfo.Initial == 0 && photoInfo.Status_id == 1) || // If photo is public and verified, show it
		(role_level == model.Role_level_Administrator) || // If user is admin, show it
		(photoInfo.Owner_id == currentUID) { // If it belongs to user, show it
		if photoInfo.Status_id == 1 {
			return true, true, nil
		}
		return true, false, nil
	}

	return false, false, nil
}

func WatermarkImagesGET(w http.ResponseWriter, r *http.Request) {
	var params = context.Get(r, "params").(httprouter.Params)
	userid := params.ByName("userid")
	pic_id := params.ByName("picid")
	user_id, _ := strconv.Atoi(userid)

	if allowed, mark, err := photoAccessAllowed(r, uint64(user_id), pic_id); allowed {
		buffer, err := renderImage(w, r, userid, pic_id, mark)
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}

		// Display JPEG to the screen
		w.Header().Set("Content-Type", "image/jpg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
			Error404(w, r)
			return
		}
	} else if err != sql.ErrNoRows {
		log.Println(err)
		Error500(w, r)
		return
	} else {
		//log.Println("User does not have access to the photo.")
		Error401(w, r)
		return
	}

}

func InitialPhotoGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	user_id := uint64(sess.Values["id"].(uint32))

	demo, err := model.DemographicByUserId(user_id)
	if err != sql.ErrNoRows {
		//log.Println(err)
	}

	// Force the user to enter in demographic information
	if len(demo.Gender) < 1 {
		UserInformationGET(w, r)
		return
	}

	// If the user has no photos, show this page
	// If the user has only unverified photos, show the waiting screen

	// Get the user photos
	photos, err := model.PhotosByUserId(uint64(sess.Values["id"].(uint32)))
	if err != nil {
		log.Println(err)
	}

	verified_private := false
	unverified_private := false
	//rejected_private := false
	any_private := false

	for _, v := range photos {
		if v.Initial == 1 {
			if v.Status_id == 1 {
				verified_private = true
			} else if v.Status_id == 2 {
				unverified_private = true
			} else if v.Status_id == 3 {
				//rejected_private = true
			}
			any_private = true
		}
	}

	// Redirect to profile to handle caess where all private photos are rejected
	if len(photos) < 1 || verified_private || !any_private {
		// Get the user verification code
		token_info, err := model.UserTokenByUserId(user_id)
		if err == sql.ErrNoRows {
			token_info.Token = random.Generate(6)
			token_info.User_id = uint32(user_id)
			err = model.UserTokenCreate(user_id, token_info.Token)
		} else if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}

		// Display the view
		v := view.New(r)
		v.Name = "user_step1"
		v.Vars["user_token"] = token_info.Token
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Render(w)
	} else if unverified_private {
		http.Redirect(w, r, "/profile", http.StatusFound)
	} else {
		//Error404(w, r)
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func InitialPhotoDeleteGET(w http.ResponseWriter, r *http.Request) {
	deletePhoto(w, r)
	http.Redirect(w, r, "/profile/initial", http.StatusFound)
}

func PhotoUploadGET(w http.ResponseWriter, r *http.Request) {
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

	// Only allow access to this page if verified
	if verified {
		// Display the view
		v := view.New(r)
		v.Name = "user_upload"
		v.Render(w)
	} else {
		Error404(w, r)
	}
}

// Displays the default home page
func PhotoPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the user photos
	photos, err := model.PhotosByUserId(uint64(sess.Values["id"].(uint32)))
	if err != nil {
		sess.AddFlash(view.Flash{"An error with the server occurred. Please try again later.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	// Limit the number of photos
	if len(photos) >= photoLimit {
		sess.AddFlash(view.Flash{"You can only have a max of " + fmt.Sprintf("%v", photoLimit) + " photos. Delete old photos and then try again.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	// File upload max size
	if r.ContentLength > 1000000*5 {
		sess.AddFlash(view.Flash{"Photo size is too large. Make sure it is under 5MB.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	// Get the form photo
	file, _, err := r.FormFile("photo")

	if err != nil {
		sess.AddFlash(view.Flash{"Photo is missing.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	defer file.Close()

	ok, filetype, _ := isSupported(file)

	// Is file supported
	if !ok {
		sess.AddFlash(view.Flash{"Photo type is not supported. Try to upload a JPG, GIF, or PNG.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	// Get the photo size
	photo_info, err := photo.ImageDimensions(file)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"Could not read the photo dimensions.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return

	}

	// OKCupid 400 x 400
	// ChristianMingle ?

	if photo_info.Width < 300 || photo_info.Height < 300 {
		sess.AddFlash(view.Flash{"Photo is too small. It must be atleast 300x300 pixels.", view.FlashError})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	user_id := fmt.Sprint(sess.Values["id"])
	folder := photoPath + user_id

	// If folder does not exists
	if !fs.FolderExists(folder) {
		err = os.Mkdir(folder, 0777)
		if err != nil {
			log.Println("Unable to create the folder for writing. Check your write access privilege.", err)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
			Index(w, r)
			return
		}
	}

	filename := time.Now().Format("20060102150405")

	finalOut := folder + "/" + filename + ".jpg"

	if filetype == "image/gif" {
		img, err := photo.GIFToImage(file)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
			Index(w, r)
			return
		}

		err = photo.ImageToJPGFile(img, finalOut)
	} else if filetype == "image/png" {
		img, err := photo.PNGToImage(file)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
			Index(w, r)
			return
		}
		err = photo.ImageToJPGFile(img, finalOut)
	} else {
		err = photo.JPGToFile(file, finalOut)
	}

	if err != nil {
		log.Println("Error uploading file:", err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
	} else {
		uid, err := strconv.ParseUint(user_id, 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
			Index(w, r)
			return
		}

		initial := false

		if strings.Contains(r.URL.Path, "initial") {
			initial = true
		}

		err = model.PhotoCreate(uid, filename, initial)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
			Index(w, r)
			return
		}

		err = photo.FixRotation(finalOut)
		if err != nil {
			//log.Println("No rotation:", err, finalOut)
		} else {
			//log.Println("Rotation success", finalOut)
		}

		po, err := pushover.New()
		if err == pushover.ErrPushoverDisabled {
			// Nothing
		} else if err != nil {
			log.Println(err)
		} else {
			err = po.Message("User " + user_id + " added a new photo for verification. You can approve the photo here:\nhttps://verified.ninja/admin/user/" + user_id)
			if err != nil {
				log.Println(err)
			}
		}

		//log.Println("File uploaded successfully:", finalOut)
		sess.AddFlash(view.Flash{"Photo uploaded successfully.", view.FlashSuccess})
	}

	sess.Save(r, w)
	Index(w, r)
	return
}

func deletePhoto(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	//userid := params.ByName("userid")
	//userid := uint64(site_idInt)
	userid := uint64(sess.Values["id"].(uint32))

	picid := params.ByName("picid")

	err := model.PhotoDelete(userid, picid)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {

		/*err = os.Remove(photoPath + fmt.Sprintf("%v", userid) + "/" + picid + ".jpg")
		if err != nil {
			log.Println(err)
		}*/

		sess.AddFlash(view.Flash{"Photo removed!", view.FlashSuccess})
		sess.Save(r, w)
	}
}

func PhotoDeleteGET(w http.ResponseWriter, r *http.Request) {
	deletePhoto(w, r)

	// Display the view
	v := view.New(r)
	v.SendFlashes(w)
}

func PhotoDownloadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params = context.Get(r, "params").(httprouter.Params)
	//userid := params.ByName("userid")
	pic_id := params.ByName("picid")
	//user_id, _ := strconv.Atoi(userid)

	user_id := uint64(sess.Values["id"].(uint32))
	userid := strconv.Itoa(int(user_id))

	if allowed, mark, err := photoAccessAllowed(r, user_id, pic_id); allowed {
		buffer, err := renderImage(w, r, userid, pic_id, mark)
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}
		// Force download
		w.Header().Set("Content-Disposition", `attachment; filename="`+pic_id+`"`)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		http.ServeContent(w, r, pic_id, time.Now(), bytes.NewReader(buffer.Bytes()))
	} else if err != sql.ErrNoRows {
		log.Println(err)
		Error500(w, r)
		return
	} else {
		//log.Println("User does not have access to the photo.")
		Error401(w, r)
		return
	}
}

/*******************************************************************************
Helpers
*******************************************************************************/

func isSupported(file multipart.File) (bool, string, error) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)

	// Reset the file
	file.Seek(0, 0)

	if err != nil {
		return false, "", err
	}

	filetype := http.DetectContentType(buff)

	for _, i := range supportedFileTypes {
		if i == filetype {
			return true, filetype, nil
		}
	}

	log.Println("File not supported: " + filetype)
	return false, filetype, nil
}
