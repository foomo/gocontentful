package erm

import (
	"os"
	"path"
	"text/template"
)

// GenerateVo generates the value objects for the space
func GenerateVo(conf SpaceConf) (err error) {
	// VO base
	tmpl, err := template.New(VoBase + TplExt).Funcs(conf.FuncMap).ParseFiles(path.Dir(conf.Filename) + TplDir + VoBase + TplExt)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Dir(conf.Filename) + OutDir + conf.PackageName + "/" + VoBase + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}

	// VO of space types

	tmpl, err = template.New(VoSpace + TplExt).Funcs(conf.FuncMap).ParseFiles(path.Dir(conf.Filename) + TplDir + VoSpace + TplExt)
	if err != nil {
		panic(err)
	}

	f, err = os.Create(path.Dir(conf.Filename) + OutDir + conf.PackageName + "/" + VoSpace + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}

	return
}
