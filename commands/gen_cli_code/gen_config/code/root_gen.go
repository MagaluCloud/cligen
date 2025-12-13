package code

import (
	_ "embed"
	"html/template"
)

//go:embed root_gen.template
var rootGenTemplate string

var rootGenTmpl *template.Template

func init() {
	var err error
	rootGenTmpl, err = template.New("rootgen").Parse(rootGenTemplate)
	if err != nil {
		panic(err)
	}
}
