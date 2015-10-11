package plugin

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"

	"github.com/verifiedninja/webapp/shared/recaptcha"
	"github.com/verifiedninja/webapp/shared/view"
)

// TemplateFuncMap returns a map of functions that are usable in templates
func TemplateFuncMap(v view.View) template.FuncMap {
	f := make(template.FuncMap)

	f["JS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("<!-- JS Error: " + s + " -->")
		}

		return template.HTML(`<script type="text/javascript" src="` + path + `"></script>`)
	}

	f["CSS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + s + " -->")
		}

		return template.HTML(`<link rel="stylesheet" type="text/css" href="` + path + `" />`)
	}

	f["LINK"] = func(path, name string) template.HTML {
		return template.HTML(`<a href="` + v.PrependBaseURI(path) + `">` + name + `</a>`)
	}

	f["SITEKEY"] = func() template.HTML {
		if recaptcha.ReadConfig().Enabled {
			return template.HTML(recaptcha.ReadConfig().SiteKey)
		}

		return template.HTML("")
	}

	f["NOESCAPE"] = func(name string) template.HTML {
		return template.HTML(name)
	}

	f["RANDIMG"] = func() template.HTML {
		num := rand.Intn(11)
		return template.HTML(fmt.Sprintf("%v", num))
	}

	f["RANDIMGSLIDER"] = func() template.HTML {
		num := rand.Intn(4)
		return template.HTML(fmt.Sprintf("%v", num))
	}

	return f
}
