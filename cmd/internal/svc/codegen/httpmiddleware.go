package codegen

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var httpMwTmpl = templates.EditableHeaderTmpl + `package httpsrv`

// GenHttpMiddleware generates http middleware file
func GenHttpMiddleware(dir string) {
	var (
		err     error
		mwfile  string
		f       *os.File
		tpl     *template.Template
		httpDir string
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	mwfile = filepath.Join(httpDir, "middleware.go")
	if _, err = os.Stat(mwfile); os.IsNotExist(err) {
		if f, err = os.Create(mwfile); err != nil {
			panic(err)
		}
		defer f.Close()
		tpl, _ = template.New("middleware.go.tmpl").Parse(httpMwTmpl)
		_ = tpl.Execute(f, struct {
			Version string
		}{
			Version: version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", mwfile)
	}
}
