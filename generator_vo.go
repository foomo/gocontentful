package erm

import (
	"os"
	"path"
	"text/template"
)

// GenerateVo generates the value objects for the space
func GenerateVo(conf SpaceConf) (err error) {
	tmpl, err := template.New(VoBase + TplExt).ParseFiles(path.Dir(conf.Filename) + TplDir + VoBase + TplExt)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Dir(conf.Filename) + OutDir + VoBase + GoExt)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}
	return
}
