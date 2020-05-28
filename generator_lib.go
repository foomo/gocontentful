package erm

import (
	"os"
	"path"
	"strings"
	"text/template"
)

// GenerateLib generates the API to access the value objects
func GenerateLib(conf SpaceConf) (err error) {
	tmpl, err := template.New(VoLib + TplExt).Funcs(conf.FuncMap).ParseFiles(path.Dir(conf.Filename) + TplDir + VoLib + TplExt)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(path.Dir(conf.Filename) + OutDir + conf.PackageName + "/" + VoLib + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}
	for _, contentType := range conf.ContentTypes {
		tmpl, err := template.New(VoLibContentType + TplExt).Funcs(conf.FuncMap).ParseFiles(path.Dir(conf.Filename) + TplDir + VoLibContentType + TplExt)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(path.Dir(conf.Filename) + OutDir + conf.PackageName + "/" + VoLib + "_" + strings.ToLower(contentType.Sys.ID) + GoExt)
		if err != nil {
			panic(err)
		}
		conf.ContentType = contentType
		err = tmpl.Execute(f, conf)
		if err != nil {
			panic(err)
		}
	}
	return
}
