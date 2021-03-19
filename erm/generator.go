package erm

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/foomo/gocontentful/erm/templates"
	"golang.org/x/tools/imports"
)

func formatAndFixImports(filename string) error {
	sourceBytes, errReadFile := ioutil.ReadFile(filename)
	if errReadFile != nil {
		return errReadFile
	}
	formattedSource, errFormat := format.Source(sourceBytes)
	if errFormat != nil {
		return errFormat
	}
	finalSource, errProcess := imports.Process(filename, formattedSource, nil)
	if errProcess != nil {
		return errProcess
	}
	return ioutil.WriteFile(filename, finalSource, 0644)
}

func generate(filename string, tpl []byte, conf spaceConf) error {
	fmt.Println("Processing", filename)
	tmpl, err := template.New("generate-" + filename).Funcs(conf.FuncMap).Parse(string(tpl))
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, conf)
	if err != nil {
		return err
	}
	errFormatAndFix := formatAndFixImports(filename)
	if errFormatAndFix != nil {
		return errFormatAndFix
	}

	return nil
}

// generateCode generates API to and value objects for the space
func generateCode(conf spaceConf) (err error) {
	for file, tpl := range map[string][]byte{
		filepath.Join(conf.PackageDir, "gocontentfulvobase"+goExt): templates.TemplateVoBase,
		filepath.Join(conf.PackageDir, "gocontentfulvo"+goExt):     templates.TemplateVo,
		filepath.Join(conf.PackageDir, "gocontentfulvolib"+goExt):  templates.TemplateVoLib,
	} {
		errGenerate := generate(file, tpl, conf)
		if errGenerate != nil {
			return errGenerate
		}
	}
	for _, contentType := range conf.ContentTypes {
		conf.ContentType = contentType
		errGenerate := generate(
			filepath.Join(conf.PackageDir, "gocontentfulvolib"+strings.ToLower(contentType.Sys.ID)+goExt),
			templates.TemplateVoLibContentType,
			conf,
		)
		if errGenerate != nil {
			return errGenerate
		}
	}
	return
}
