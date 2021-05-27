package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func ReadFileToAst(path string) (f *ast.File, err error) {
	src := ReadFile(path)
	fset := token.NewFileSet()
	f, err = parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	ast.Print(fset, f)
	return f, nil
}

func GetCommentAndName(genDecl *ast.GenDecl) (string, string) {
	comment, tableName := "", ""
	for _, field := range genDecl.Doc.List { //基于机结构体的注释
		text := GetStructComment(field.Text)
		if len(text) < 2 {
			continue
		}
		switch strings.ToLower(text[0]) {
		case "comment":
			comment = "'" + text[1] + "'"
		case "name":
			tableName = text[1]
		}
	}
	return comment, tableName
}

func GeneratorSql() {
	f, _ := ReadFileToAst("user.go")
	data := new(SqlData)
	for _, node := range f.Decls {
		genDecl := node.(*ast.GenDecl)
		desc := new(SqlDataChild)
		desc.Comment, desc.TableName = GetCommentAndName(genDecl) //基于机结构体的注释

		for _, field := range genDecl.Specs {
			switch t := field.(type) {
			case *ast.TypeSpec: //表名
				if desc.TableName == "" {
					desc.TableName = t.Name.Name
				}
				switch tt := t.Type.(type) {
				case *ast.StructType:
					desc = StructToSql(tt, desc)
				}
			}
		}
		data.DescList = append(data.DescList, desc)
	}
	CreateSql(data)
}

func StructToSql(tt *ast.StructType, desc *SqlDataChild) *SqlDataChild {

	for i, field := range tt.Fields.List {
		fieldName := field.Names[0].Name
		switch fieldName {
		case "Id": //默认Id是主键
			if desc.PrimaryKey == "" {
				desc.PrimaryKey = fieldName
			}
		}
		//判断tag是否配置primaryKey
		tagMap := FiledToMap(field.Tag.Value)
		if _, ok := tagMap["primaryKey"]; ok {
			desc.PrimaryKey = fieldName
			delete(tagMap, "primaryKey")
		}
		//是否配置了column
		if columnName, ok := tagMap["column"]; ok {
			if desc.PrimaryKey == fieldName {
				desc.PrimaryKey = columnName
			}
			fieldName = columnName
			delete(tagMap, "column")
		}

		//是否配置了size
		fieldTYpe := GetColumnType(field.Type.(*ast.Ident).Name)
		if size, ok := tagMap["size"]; ok {
			if fieldTYpe == "varchar(255)" {
				fieldTYpe = "varchar(" + size + ")"
			} else {
				fieldTYpe += "(" + size + ")"
			}
			delete(tagMap, "size")
		}
		desc.List = append(desc.List, &SqlDesc{
			Index:   i,
			Name:    fieldName,
			Type:    fieldTYpe,
			Tag:     GetFieldTag(tagMap),
			Comment: GetColumnComment(field.Comment.Text()),
		})
	}
	return desc
}

//data.DescList = append(data.DescList, desc)
func GetStructComment(str string) []string {
	text := ReplaceStr(str, "//@", "")
	text = strings.TrimSpace(text)
	return strings.Split(text, " ")
}

func FiledToMap(srcTag string) (tm map[string]string) {
	tm = make(map[string]string, 0)
	srcTag = CleanQuote(srcTag)
	srcTag = strings.TrimPrefix(srcTag, "db:")
	srcTag = CleanDoubleQuotes(srcTag)
	for _, s := range strings.Split(srcTag, ";") {
		preS := strings.Split(s, ":")
		if len(preS) == 2 {
			tm[preS[0]] = preS[1]
		} else {
			tm[s] = s
		}
	}
	return
}
func GetFieldTag(srcTag map[string]string) string {
	targetTag := ""
	for _, s := range srcTag {
		targetTag += s + " "
	}
	return targetTag
}

func GetColumnComment(comment string) string {
	return "COMMENT" + "'" + strings.TrimSpace(comment) + "'"
}

func GetColumnType(name string) string {
	return GoTypeToMysqlType[name]
}

func CreateSql(data *SqlData) {
	// 写入markdown
	dir, _ := os.Getwd()
	file := dir + "/" + "struct.sql"
	CreateFileIfHasDel(file)
	tplByte, err := ioutil.ReadFile(TPL_SQL)
	if err != nil {
		return
	}
	// 解析
	content := bytes.NewBuffer([]byte{})
	tpl, err := template.New("sql").Parse(string(tplByte))
	err = tpl.Execute(content, data)
	if err != nil {
		return
	}
	// 表信息写入文件
	err = WriteAppendFile(file, content.String())
	if err != nil {
		return
	}
	return

}
