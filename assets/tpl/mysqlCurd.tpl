package model

import (
"database/sql"
)

type {{.StructTableName}}Model struct {
    Entry *{{.StructTableName}}
    columns []string
    columnMap map[string]interface{}

    pageSize  int
    limit     string
    where    []string
    or       []string
    joins   string
    set      []string
    val []interface{}
    DB *sql.DB
    Tx *sql.Tx
}

func New{{.StructTableName}}(db *sql.DB) *{{.StructTableName}}Model {
     return &{{.StructTableName}}Model{
                DB: db,
                Entry:&{{.StructTableName}}{},
                columnMap: make(map[string]interface{}, 0),
      }
}

func New{{.StructTableName}}WithEntry(db *sql.DB, entry *{{.StructTableName}}) *{{.StructTableName}}Model  {
	return &{{.StructTableName}}Model {
		DB:        db,
		Entry:     entry,
		columnMap: make(map[string]interface{}, 0),
	}
}

func New{{.StructTableName}}Tx(tx *sql.Tx) *{{.StructTableName}}Model {
	return &{{.StructTableName}}Model{
		Tx: tx,
		 Entry:&{{.StructTableName}}{},
	}
}

// 获取所有的表字段
func (m *{{.StructTableName}}Model) getColumns() string {
    return " {{.AllFieldList}} "
}
//设置查询字段
func (m *{{.StructTableName}}Model) Columns(column ...string) *{{.StructTableName}}Model {
	for _, c := range column {
        m.columns = append(m.columns, c)

         {{range $k, $v :=.FieldsInfo}}
         if c=="{{$v.Name}}"{
            m.columnMap[c] = &m.Entry.{{$v.HumpName}}
         }
         {{end}}
    }
	return m
}
//设置查询字段
func (m *{{.StructTableName}}Model) Where(field string, val interface{}) *{{.StructTableName}}Model {
	m.where = append(m.where, field)
    m.val = append(m.val, val)
	return m
}

//设置查询字段
func (m *{{.StructTableName}}Model) And(field string, val interface{}) *{{.StructTableName}}Model {
	m.where = append(m.where, field)
    	m.val = append(m.val, val)
	return m
}

//设置Or
func (m *{{.StructTableName}}Model) Or(field string, val interface{})*{{.StructTableName}}Model  {
	m.or = append(m.or, field)
	m.val = append(m.val, val)
	return m
}

func (m *{{.StructTableName}}Model) SetPageSize(size int) *{{.StructTableName}}Model {
	m.pageSize = size
	return m
}
func (m *{{.StructTableName}}Model) Page(page int) *{{.StructTableName}}Model {
	//(page-1)*pageSize, pageSize
	offset := (page - 1) * m.pageSize
	m.limit = " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(m.pageSize)
	return m
}
//设置联合查询的条件
func (m *{{.StructTableName}}Model) Joins(join string) *{{.StructTableName}}Model {
	m.joins = " " + join + " "
	return m
}



//生成set、get方法
 {{range $key, $item := .FieldsInfo}}
 func(m *{{$.StructTableName}}Model)Set{{$item.HumpName}}({{$item.Name}} {{$item.FieldType}})*{{$.StructTableName}}Model{
    m.Entry.{{$item.HumpName}} = {{$item.Name}}
    return m
 }
 func (m *{{$.StructTableName}}Model)Get{{$item.HumpName}}(){{$item.FieldType}}{
 	return m.Entry.{{$item.HumpName}}
 }
 {{end}}



// 获取多行数据.
func (m *{{.StructTableName}}Model) getRows(ctx context.Context,sqlTxt string, params ...interface{}) (rowsResult []*{{.StructTableName}}, err error) {
    query, err := m.DB.QueryContext(ctx,sqlTxt, params...)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer query.Close()

    columns, _ := query.Columns()
    dest := make([]interface{}, 0, len(columns))
    for _, v := range columns {
    	dest = append(dest, m.columnMap[v])
    }


    for query.Next() {
        row := {{.StructTableName}}{}
        err = query.Scan(dest...)
        if err != nil {
           fmt.Println(err)
            return
        }

        {{range .FieldsInfo}}row.{{.HumpName}}=m.Entry.{{.HumpName}}// {{.Comment}}
        {{end}}
        rowsResult = append(rowsResult, &row)
    }
    return
}

//获取单行
func (m *{{.StructTableName}}Model) getRow(ctx context.Context,db *sql.DB,sqlText string, params ...interface{}) (rowResult *{{.StructTableName}}, err error) {
    query := m.DB.QueryRowContext(ctx,sqlText, params...)
        row := {{.StructTableName}}{}
        err = query.Scan(
       {{range .FieldsInfo}}&row.{{.HumpName}},// {{.Comment}}
        {{end}})
        if err != nil {
        fmt.Println(err)
        return
        }
        rowResult = &row
        return
}

// 保存数据
func (m *{{.StructTableName}}Model) Save(ctx context.Context,sqlTxt string, value ...interface{}) (b bool, err error) {
    stmt, err := m.DB.PrepareContext(ctx, sqlTxt)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer stmt.Close()
    result, err := stmt.Exec(value...)
    if err != nil {
        fmt.Println(err)
        return
    }
    var affectCount int64
    affectCount, err = result.RowsAffected()
    if err != nil {
        fmt.Println(err)
        return
    }
    b = affectCount > 0
    return
}

// _更新数据Tx
func (m *{{.StructTableName}}Model) SaveTx(ctx context.Context,tx *sql.Tx,sqlTxt string, value ...interface{}) (b bool, err error) {
    stmt, err := tx.PrepareContext(ctx,sqlTxt)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer stmt.Close()
    result, err := stmt.Exec(value...)
    if err != nil {
        _=tx.Rollback()
        fmt.Println(err)
        return
    }
    var affectCount int64
    affectCount, err = result.RowsAffected()
    if err != nil {
        fmt.Println(err)
        return
    }
    b = affectCount > 0
    return
}

// 新增信息
func (m *{{.StructTableName}}Model) Create() (lastId int64, err error) {
    const sqlText = "INSERT INTO " + "{{.UpperTableName}}" + " ({{.InsertFieldList}}) VALUES ({{.InsertMark}})"
    stmt, err := m.DB.Prepare(sqlText)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer stmt.Close()
    result, err := stmt.Exec(
    {{range .InsertInfo}}m.Entry.{{.HumpName}},// {{.Comment}}
    {{end}})
    if err != nil {
        fmt.Println(err)
        return
    }
    lastId, err = result.LastInsertId()
    if err != nil {
        fmt.Println(err)
        return
    }
    return
}

// 获取单行数据
func (m *{{.StructTableName}}Model) Find(ctx context.Context,id ...int64) (result []*{{.StructTableName}}, err error) {
	columns := strings.Join(m.columns,",")
	if columns == "" {
		columns = m.getColumns()
	}

	whereField:= m.getWhere()
	//if whereField == "" {
		//whereField = " where Id = ?"
		//val = append(val, id)
	//}

	or := m.getOr()

	limit := m.limit
    if limit == "" {
    	limit = " LIMIT 1"
    }
    join := m.joins
    if join != "" {
    	limit = ""
    }
    sqlText := "SELECT " + columns + " FROM " + "user" + whereField + or + join + limit

	result, err = m.getRows(ctx, sqlText, m.val...)
	return
}

//设置更新字段
func (m *{{.StructTableName}}Model) Set(fields string,val interface{}) *{{.StructTableName}}Model {
	m.set = append(m.set, tools.AddQuote(fields)+"=?")
    m.val = append(m.val, val)
	return m
}

// 更新数据
func (m *{{.StructTableName}}Model) Update(ctx context.Context) (b bool, err error) {
    set := strings.Join(m.set, ",")
    whereField := m.getWhere()

    sqlText := "UPDATE " + " user " + " SET " + set + whereField

    return m.Save(ctx,sqlText, m.val...)
}



func (m *{{.StructTableName}}Model) getWhere() string {
    if len(m.where)==0{
		return ""
	}
   return " where " + strings.Join(m.where, " and ")
}

func (m *{{.StructTableName}}Model)getOr() string {
    if len(m.or) == 0 {
		return ""
	}
	return " or " + strings.Join(m.or, " or ")
}
