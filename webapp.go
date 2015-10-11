package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/verifiedninja/webapp/route"
	"github.com/verifiedninja/webapp/shared/database"
	"github.com/verifiedninja/webapp/shared/email"
	"github.com/verifiedninja/webapp/shared/jsonconfig"
	"github.com/verifiedninja/webapp/shared/pushover"
	"github.com/verifiedninja/webapp/shared/recaptcha"
	"github.com/verifiedninja/webapp/shared/server"
	"github.com/verifiedninja/webapp/shared/session"
	"github.com/verifiedninja/webapp/shared/view"
	"github.com/verifiedninja/webapp/shared/view/plugin"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Configure the SMTP server
	email.Configure(config.Email)

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Configure Pushover
	pushover.Configure(config.Pushover)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(plugin.TemplateFuncMap(config.View))

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.Databases      `json:"Database"`
	Email     email.SMTPInfo          `json:"Email"`
	Recaptcha recaptcha.RecaptchaInfo `json:"Recaptcha"`
	Pushover  pushover.PushoverInfo   `json:"Pushover"`
	Server    server.Server           `json:"Server"`
	Session   session.Session         `json:"Session"`
	Template  view.Template           `json:"Template"`
	View      view.View               `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
