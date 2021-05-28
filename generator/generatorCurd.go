package generator

import (
	"bytes"
	"github.com/1975210542/generatorTools/config"
	"github.com/1975210542/generatorTools/entry"
	"github.com/1975210542/generatorTools/tools"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func GeneratorCurd(scanPath, outputPath string) {
	fset := token.NewFileSet()
	fs, _ := parser.ParseDir(fset, scanPath, nil, parser.ParseComments)
	for _, ff := range fs {
		for _, f := range ff.Files {
			generatorOne(f, outputPath)
		}
	}
}

func generatorOne(f *ast.File, outputPath string) {
	data := new(entry.CurdSqlInfo)
	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.TypeSpec:
			data.StructTableName = t.Name.Name
			data.TableName = tools.Lcfirst(t.Name.Name)
			data.UpperTableName = strings.ToLower(t.Name.Name)
		case *ast.StructType:
			allFields := []string{}
			insertFieldList := []string{}
			insertMark := []string{}
			for _, field := range t.Fields.List {
				fieldName := field.Names[0].Name

				switch fieldName {
				case "Id": //默认Id是主键
					if data.PrimaryKey == "" {
						data.PrimaryKey = fieldName
					}
				}
				//判断tag是否配置primaryKey
				if field.Tag == nil {
					continue
				}
				tagMap := filedToMap(field.Tag.Value)
				if _, ok := tagMap["primaryKey"]; ok {
					data.PrimaryKey = fieldName
				}
				allFields = append(allFields, tools.AddQuote(strings.ToLower(fieldName)))

				sqlFieldInfo := &entry.SqlFieldInfo{
					HumpName:  field.Names[0].Name,
					Name:      tools.Lcfirst(field.Names[0].Name),
					FieldType: field.Type.(*ast.Ident).Name,
					NullType:  entry.GoTypeToMysqlNullType[field.Type.(*ast.Ident).Name],
					Comment:   field.Comment.Text(),
				}

				if fieldName != "Id" {
					insertFieldList = append(insertFieldList, tools.AddQuote(strings.ToLower(fieldName)))
					insertMark = append(insertMark, "?")
					//插入的字段
					data.InsertInfo = append(data.InsertInfo, sqlFieldInfo)
				}
				//所有的字段
				data.FieldsInfo = append(data.FieldsInfo, sqlFieldInfo)
			}

			data.AllFieldList = strings.Join(allFields, ",")
			data.InsertFieldList = strings.Join(insertFieldList, ",")
			data.InsertMark = strings.Join(insertMark, ",")
		}
		return true
	})
	createCURD(data, outputPath)

}

func createCURD(data *entry.CurdSqlInfo, outputPath string) {
	// 写入markdown
	dir, _ := os.Getwd()
	file := dir + "/"+outputPath + "/" + data.TableName + "Model.go"
	tools.CreateFileIfHasDel(file)
	tplByte, err := ioutil.ReadFile(config.TPL_My_CURD)
	if err != nil {
		return
	}
	// 解析
	content := bytes.NewBuffer([]byte{})
	tpl, err := template.New("curd").Parse(string(tplByte))
	err = tpl.Execute(content, data)
	if err != nil {
		return
	}
	// 表信息写入文件
	err = tools.WriteAppendFile(file, content.String())
	if err != nil {
		return
	}

	tools.Gofmt(file)
	return

}
