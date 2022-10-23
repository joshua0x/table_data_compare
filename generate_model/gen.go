package generate_model

import (
	"fmt"
	"github.com/joshua0x/table_data_compare/config"
	parser "github.com/joshua0x/table_data_compare/parser_from_cascax"
	"github.com/pkg/errors"
	"html/template"
	"os"
	"strings"
)

type TableModel struct {
	TableStructName string
	TableStructDef template.HTML
	TableName string
	IdObjName string
	PkColumnName string
	Tag template.HTML
	ImportPath []string
}


func convert(src parser.ModelCodes,tn string) TableModel {
	tm := TableModel{}
	tm.TableStructDef = template.HTML(strings.Join(src.StructCode,"\n"))
	tm.TableStructName = src.TableCamelName
	tm.TableName = tn
	tm.PkColumnName = src.PkColName
	tm.IdObjName = "IdObj" + src.TableCamelName
	tm.ImportPath = src.ImportPath
	//
	builder := strings.Builder{}
	builder.WriteString("gorm")
	builder.WriteString(`:"`)
	builder.WriteString(src.PkColName)
	builder.WriteString(`"`)

	tm.Tag = template.HTML(builder.String())
	return tm
}

func Gen(cfg *config.DbConfig) {
	//??,list-all-tablenames??,
	template.HTMLEscaper()
	temp,err := template.New("x").Parse(templateStr)
	if err != nil {
		panic(err)
	}
	//

	for _,table := range cfg.TableList {
		res, err := parser.ParseSqlFromDB(cfg.HostA, table)
		if err != nil {
			panic(err)
		}
		if res.TableCamelName == "" {
			panic(errors.New("got-errors"))
		}
		//how to convered
		fname := "./tab_models/"+table + ".go"
		fobj,err := os.Create(fname)
		if err != nil {
			panic(err)
		}

		err2 := temp.Execute(fobj,convert(res,table))
		if err2 != nil {
			panic(err2)
		}
		fmt.Printf("finished generate go-struct for table:%v in %v\n",table,fname)

	}
}