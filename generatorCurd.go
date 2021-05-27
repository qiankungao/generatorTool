package main

import (
	"bytes"
	"go/ast"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type CurdSqlInfo struct {
	TableName           string              // 表名
	PrimaryKey          string              // 主键字段
	PrimaryType         string              // 主键类型
	StructTableName     string              // 结构表名称
	NullStructTableName string              // 判断为空的表名
	PkgEntity           string              // 实体空间名称
	PkgTable            string              // 表的空间名称
	UpperTableName      string              // 大写的表名
	AllFieldList        string              // 所有字段列表,如: id,name
	InsertFieldList     string              // 插入字段列表,如:id,name
	InsertMark          string              // 插入字段使用多少个?,如: ?,?,?
	UpdateFieldList     string              // 更新字段列表
	SecondField         string              // 存放第二个字段
	UpdateListField     []string            // 更新字段列表
	FieldsInfo          []*SqlFieldInfo     // 字段信息
	NullFieldsInfo      []*NullSqlFieldInfo // 判断为空时
	InsertInfo          []*SqlFieldInfo
}

func GenCurd() {
	/*
		const sqlText = "INSERT INTO " + config.TABLE_BORDER + " (`playerId`,`stageInfo`,`freeTimes`) VALUES (?,?,?)"
	*/

	f, _ := ReadFileToAst("user.go")
	data := new(CurdSqlInfo)
	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.TypeSpec:
			data.StructTableName = t.Name.Name
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
				tagMap := FiledToMap(field.Tag.Value)
				if _, ok := tagMap["primaryKey"]; ok {
					data.PrimaryKey = fieldName
				}
				allFields = append(allFields, AddQuote(strings.ToLower(fieldName)))

				sqlFieldInfo := &SqlFieldInfo{
					HumpName:  field.Names[0].Name,
					Name:      Lcfirst(field.Names[0].Name),
					FieldType: field.Type.(*ast.Ident).Name,
					NullType:  GoTypeToMysqlNullType[field.Type.(*ast.Ident).Name],
					Comment:   field.Comment.Text(),
				}

				if fieldName != "Id" {
					insertFieldList = append(insertFieldList, AddQuote(strings.ToLower(fieldName)))
					insertMark = append(insertMark, "?")
					//插入的字段
					data.InsertInfo = append(data.InsertInfo, sqlFieldInfo)
				}
				//所有的字段
				data.FieldsInfo = append(data.FieldsInfo, sqlFieldInfo)

				data.NullFieldsInfo = append(data.NullFieldsInfo, &NullSqlFieldInfo{
					HumpName: field.Names[0].Name,
					GoType:   field.Type.(*ast.Ident).Name,
					Comment:  field.Comment.Text(),
				})

			}

			data.AllFieldList = strings.Join(allFields, ",")
			data.InsertFieldList = strings.Join(insertFieldList, ",")
			data.InsertMark = strings.Join(insertMark, ",")
		}
		return true
	})

	CreateCURD(data)

}

func CreateCURD(data *CurdSqlInfo) {
	// 写入markdown
	dir, _ := os.Getwd()
	file := dir + "/" + "struct.go"
	CreateFileIfHasDel(file)
	tplByte, err := ioutil.ReadFile(TPL_My_CURD)
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
	err = WriteAppendFile(file, content.String())
	if err != nil {
		return
	}

	Gofmt(file)
	return

}
