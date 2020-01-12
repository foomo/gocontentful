package erm

import (
	"os"
	"path"
	"text/template"
)

// GenerateLib generates the API to access the value objects
func GenerateLib(conf SpaceConf) (err error) {
	tmpl, err := template.New(VoLib + TplExt).Funcs(conf.FuncMap).ParseFiles(path.Dir(conf.Filename) + TplDir + VoLib + TplExt)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Dir(conf.Filename) + OutDir + VoLib + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}
	return
}
