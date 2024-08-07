package codegen

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var dbTmpl = templates.EditableHeaderTmpl + `package db

import (
	"{{.ConfigPackage}}"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewDb(conf config.DbConfig) (*sqlx.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += "&loc=Asia%2FShanghai&parseTime=True"

	db, err := sqlx.Connect(conf.Driver, conn)
	if err != nil {
		return nil, errors.Wrap(err, "database connection failed")
	}
	db.MapperFunc(strcase.ToSnake)
	return db, nil
}
`

var MkdirAll = os.MkdirAll
var Open = os.Open
var Create = os.Create
var Stat = os.Stat

// GenDb generates db connection code
func GenDb(dir string) {
	var (
		err    error
		dbfile string
		f      *os.File
		tpl    *template.Template
		dbDir  string
	)
	dbDir = filepath.Join(dir, "db")
	if err = MkdirAll(dbDir, os.ModePerm); err != nil {
		panic(err)
	}

	dbfile = filepath.Join(dbDir, "db.go")
	if _, err = Stat(dbfile); os.IsNotExist(err) {
		cfgPkg := astutils.GetPkgPath(filepath.Join(dir, "config"))
		if f, err = Create(dbfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("db.go.tmpl").Parse(dbTmpl); err != nil {
			panic(err)
		}
		_ = tpl.Execute(f, struct {
			ConfigPackage string
			Version       string
		}{
			ConfigPackage: cfgPkg,
			Version:       version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", dbfile)
	}
}
