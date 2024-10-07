package codegen

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

func genMainModule(dir string) {
	var (
		err      error
		mainfile string
		f        *os.File
		tpl      *template.Template
		cmdDir   string
		buf      bytes.Buffer
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}
	mainfile = filepath.Join(cmdDir, "main.go")
	if f, err = Create(mainfile); err != nil {
		panic(err)
	}
	defer f.Close()

	if tpl, err = template.New(templates.MainModuleTmpl).Parse(templates.MainModuleTmpl); err != nil {
		panic(err)
	}
	pluginPkg := astutils.GetPkgPath(filepath.Join(dir, "plugin"))
	if err = tpl.Execute(&buf, struct {
		PluginPackage string
		Version       string
	}{
		PluginPackage: pluginPkg,
		Version:       version.Release,
	}); err != nil {
		panic(err)
	}
	astutils.FixImport([]byte(strings.TrimSpace(buf.String())), mainfile)
}
