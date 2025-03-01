package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/table"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var idaoTmpl = templates.EditableHeaderTmpl + `package dao

import (
	"context"
	"github.com/unionj-cloud/toolkit/sqlext/query"
	"{{.EntityPackage}}"
)

type I{{.EntityName}}Dao interface {
	// single table CRUD operations
	Insert(ctx context.Context, data *entity.{{.EntityName}}) (int64, error)
	InsertIgnore(ctx context.Context, data *entity.{{.EntityName}}) (int64, error)
	BulkInsert(ctx context.Context, data []*entity.{{.EntityName}}) (int64, error)
	BulkInsertIgnore(ctx context.Context, data []*entity.{{.EntityName}}) (int64, error)
	Upsert(ctx context.Context, data *entity.{{.EntityName}}) (int64, error)
	BulkUpsert(ctx context.Context, data []*entity.{{.EntityName}}) (int64, error)
	BulkUpsertSelect(ctx context.Context, data []*entity.{{.EntityName}}, columns []string) (int64, error)
	UpsertNoneZero(ctx context.Context, data *entity.{{.EntityName}}) (int64, error)
	DeleteMany(ctx context.Context, where query.Where) (int64, error)
	Update(ctx context.Context, data *entity.{{.EntityName}}) (int64, error)
	UpdateNoneZero(ctx context.Context, data *entity.{{.EntityName}}) (int64, error)
	UpdateMany(ctx context.Context, data []*entity.{{.EntityName}}, where query.Where) (int64, error)
	UpdateManyNoneZero(ctx context.Context, data []*entity.{{.EntityName}}, where query.Where) (int64, error)
	Get(ctx context.Context, dest *entity.{{.EntityName}}, id {{.PkField.Type}}) error
	SelectMany(ctx context.Context, dest *[]entity.{{.EntityName}}, where query.Where) error
	CountMany(ctx context.Context, where query.Where) (int, error)
	PageMany(ctx context.Context, dest *{{.EntityName}}PageRet, page query.Page, where query.Where) error
	DeleteManySoft(ctx context.Context, where query.Where) (int64, error)

	// hooks
	BeforeSaveHook(ctx context.Context, data *entity.{{.EntityName}})
	BeforeBulkSaveHook(ctx context.Context, data []*entity.{{.EntityName}})
	AfterSaveHook(ctx context.Context, data *entity.{{.EntityName}}, lastInsertID int64, affected int64)
	AfterBulkSaveHook(ctx context.Context, data []*entity.{{.EntityName}}, lastInsertID int64, affected int64)
	BeforeUpdateManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where)
	AfterUpdateManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where, affected int64)
	BeforeDeleteManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where)
	AfterDeleteManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where, affected int64)
	BeforeReadManyHook(ctx context.Context, page *query.Page, where *query.Where)
}`

// GenIDaoGo generates dao layer interface code
func GenIDaoGo(entityPath string, t table.Table, folder ...string) error {
	var (
		err      error
		daopath  string
		f        *os.File
		tpl      *template.Template
		df       string
		dpkg     string
		pkColumn table.Column
	)
	df = "dao"
	if len(folder) > 0 {
		df = folder[0]
	}
	daopath = filepath.Join(filepath.Dir(entityPath), df)
	_ = os.MkdirAll(daopath, os.ModePerm)
	daofile := filepath.Join(daopath, "i"+strings.ToLower(t.Meta.Name)+"dao.go")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		f, _ = os.Create(daofile)
		defer f.Close()
		dpkg = astutils.GetImportPath(entityPath)
		for _, column := range t.Columns {
			if column.Pk {
				pkColumn = column
				break
			}
		}
		tpl, _ = template.New("idao.go.tmpl").Parse(idaoTmpl)
		_ = tpl.Execute(f, struct {
			EntityName    string
			Version       string
			EntityPackage string
			PkField       astutils.FieldMeta
		}{
			EntityName:    t.Meta.Name,
			Version:       version.Release,
			EntityPackage: dpkg,
			PkField:       pkColumn.Meta,
		})
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}
