package route

import (
	"net/http"

	"github.com/verifiedninja/webapp/controller"
	"github.com/verifiedninja/webapp/route/middleware/acl"
	hr "github.com/verifiedninja/webapp/route/middleware/httprouterwrapper"
	"github.com/verifiedninja/webapp/route/middleware/logrequest"
	//"github.com/verifiedninja/webapp/route/middleware/pprofhandler"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Load the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())
}

// Load the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page and public pages
	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.Index)))
	r.GET("/about", hr.Handler(alice.
		New().
		ThenFunc(controller.AboutGET)))
	r.GET("/terms", hr.Handler(alice.
		New().
		ThenFunc(controller.TermsGET)))
	r.GET("/contact", hr.Handler(alice.
		New().
		ThenFunc(controller.ContactGET)))
	r.POST("/contact", hr.Handler(alice.
		New().
		ThenFunc(controller.ContactPOST)))
	r.GET("/verify", hr.Handler(alice.
		New().
		ThenFunc(controller.VerifyUsernameGET)))
	r.GET("/public/:site/:username", hr.Handler(alice.
		New().
		ThenFunc(controller.PublicUsernameGET)))
	r.GET("/image/:userid/:picid", hr.Handler(alice.
		New().
		ThenFunc(controller.WatermarkImagesGET)))

	// Login and logout
	r.GET("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
	r.GET("/logout", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.Logout)))

	// Register
	r.GET("/register", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterGET)))
	r.POST("/register", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterPOST)))

	// Email verification
	r.GET("/emailverification/:token", hr.Handler(alice.
		New().
		ThenFunc(controller.EmailVerificationGET)))

	// User Pages
	r.GET("/profile", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserProfileGET)))
	r.GET("/profile/initial", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.InitialPhotoGET)))
	r.POST("/profile/initial", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PhotoPOST)))
	r.GET("/profile/initial/delete/:picid", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.InitialPhotoDeleteGET)))
	r.GET("/profile/photo/upload", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PhotoUploadGET)))
	r.POST("/profile/photo/upload", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PhotoPOST)))
	r.GET("/profile/photo/delete/:picid", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PhotoDeleteGET)))
	r.GET("/profile/site", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserSiteGET)))
	r.POST("/profile/site", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserSitePOST)))
	r.GET("/profile/photo/download/:picid", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PhotoDownloadGET)))
	r.POST("/profile/information", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserInformationPOST)))
	r.GET("/profile/information", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserInformationGET)))
	r.GET("/profile/email", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserEmailGET)))
	r.POST("/profile/email", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserEmailPOST)))
	r.GET("/profile/password", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserPasswordGET)))
	r.POST("/profile/password", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UserPasswordPOST)))

	// Admin Pages
	r.GET("/admin", hr.Handler(alice.
		New(acl.AllowOnlyAdministrator).
		ThenFunc(controller.AdminGET)))
	r.GET("/admin/user/:userid", hr.Handler(alice.
		New(acl.AllowOnlyAdministrator).
		ThenFunc(controller.AdminUserGET)))
	r.GET("/admin/all", hr.Handler(alice.
		New(acl.AllowOnlyAdministrator).
		ThenFunc(controller.AdminAllGET)))
	r.GET("/admin/approve/:userid/:picid", hr.Handler(alice.
		New(acl.AllowOnlyAdministrator).
		ThenFunc(controller.AdminApproveGET)))
	r.GET("/admin/reject/:userid/:picid", hr.Handler(alice.
		New(acl.AllowOnlyAdministrator).
		ThenFunc(controller.AdminRejectGET)))
	r.GET("/admin/unverify/:userid/:picid", hr.Handler(alice.
		New(acl.AllowOnlyAdministrator).
		ThenFunc(controller.AdminUnverifyGET)))

	// API
	r.GET("/api/v1/verify/:site/:username", hr.Handler(alice.
		New().
		ThenFunc(controller.APIVerifyUserGET)))
	r.GET("/api/v1/request/token", hr.Handler(alice.
		New().
		ThenFunc(controller.APIRegisterChromeGET)))

	// Cron
	// TODO This should not be publicly accessible
	r.GET("/cron/notifyemailexpire", hr.Handler(alice.
		New().
		ThenFunc(controller.CronNotifyEmailExpireGET)))

	// Enable Pprof
	/*r.GET("/debug/pprof/*pprof", hr.Handler(alice.
	New(acl.DisallowAnon).
	ThenFunc(pprofhandler.Handler)))*/

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Log every request
	h = logrequest.Database(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
