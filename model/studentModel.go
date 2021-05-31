package model

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/1975210542/generatorTools/tools"
)

type StudentModel struct {
	Entry     *Student
	columns   []string
	columnMap map[string]interface{}

	pageSize int
	limit    string
	where    []string
	or       []string
	joins    string
	set      []string
	val      []interface{}
	DB       *sql.DB
	Tx       *sql.Tx
}

func NewStudent(db *sql.DB) *StudentModel {
	return &StudentModel{
		DB:        db,
		Entry:     &Student{},
		columnMap: make(map[string]interface{}, 0),
	}
}

func NewStudentWithEntry(db *sql.DB, entry *Student) *StudentModel {
	return &StudentModel{
		DB:        db,
		Entry:     entry,
		columnMap: make(map[string]interface{}, 0),
	}
}

func NewStudentTx(tx *sql.Tx) *StudentModel {
	return &StudentModel{
		Tx:    tx,
		Entry: &Student{},
	}
}

// 获取所有的表字段
func (m *StudentModel) getColumns() string {
	return " `id`,`name`,`age` "
}

//设置查询字段
func (m *StudentModel) Columns(column ...string) *StudentModel {
	for _, c := range column {
		m.columns = append(m.columns, c)

		if c == "id" {
			m.columnMap[c] = &m.Entry.Id
		}

		if c == "name" {
			m.columnMap[c] = &m.Entry.Name
		}

		if c == "age" {
			m.columnMap[c] = &m.Entry.Age
		}

	}
	return m
}

//设置查询字段
func (m *StudentModel) Where(field string, val interface{}) *StudentModel {
	m.where = append(m.where, field)
	m.val = append(m.val, val)
	return m
}

//设置查询字段
func (m *StudentModel) And(field string, val interface{}) *StudentModel {
	m.where = append(m.where, field)
	m.val = append(m.val, val)
	return m
}

//设置Or
func (m *StudentModel) Or(field string, val interface{}) *StudentModel {
	m.or = append(m.or, field)
	m.val = append(m.val, val)
	return m
}

func (m *StudentModel) SetPageSize(size int) *StudentModel {
	m.pageSize = size
	return m
}
func (m *StudentModel) Page(page int) *StudentModel {
	//(page-1)*pageSize, pageSize
	offset := (page - 1) * m.pageSize
	m.limit = " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(m.pageSize)
	return m
}

//设置联合查询的条件
func (m *StudentModel) Joins(join string) *StudentModel {
	m.joins = " " + join + " "
	return m
}

//生成set、get方法

func (m *StudentModel) SetId(id int32) *StudentModel {
	m.Entry.Id = id
	return m
}
func (m *StudentModel) GetId() int32 {
	return m.Entry.Id
}

func (m *StudentModel) SetName(name string) *StudentModel {
	m.Entry.Name = name
	return m
}
func (m *StudentModel) GetName() string {
	return m.Entry.Name
}

func (m *StudentModel) SetAge(age int32) *StudentModel {
	m.Entry.Age = age
	return m
}
func (m *StudentModel) GetAge() int32 {
	return m.Entry.Age
}

// 获取多行数据.
func (m *StudentModel) getRows(ctx context.Context, sqlTxt string, params ...interface{}) (rowsResult []*Student, err error) {
	query, err := m.DB.QueryContext(ctx, sqlTxt, params...)
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
		row := Student{}
		err = query.Scan(dest...)
		if err != nil {
			fmt.Println(err)
			return
		}

		row.Id = m.Entry.Id // 主键递增

		row.Name = m.Entry.Name // 姓名

		row.Age = m.Entry.Age //

		rowsResult = append(rowsResult, &row)
	}
	return
}

//获取单行
func (m *StudentModel) getRow(ctx context.Context, db *sql.DB, sqlText string, params ...interface{}) (rowResult *Student, err error) {
	query := m.DB.QueryRowContext(ctx, sqlText, params...)
	row := Student{}
	err = query.Scan(
		&row.Id, // 主键递增

		&row.Name, // 姓名

		&row.Age, //
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	rowResult = &row
	return
}

// 保存数据
func (m *StudentModel) Save(ctx context.Context, sqlTxt string, value ...interface{}) (b bool, err error) {
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
func (m *StudentModel) SaveTx(ctx context.Context, tx *sql.Tx, sqlTxt string, value ...interface{}) (b bool, err error) {
	stmt, err := tx.PrepareContext(ctx, sqlTxt)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(value...)
	if err != nil {
		_ = tx.Rollback()
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
func (m *StudentModel) Create() (lastId int64, err error) {
	const sqlText = "INSERT INTO " + "student" + " (`name`,`age`) VALUES (?,?)"
	stmt, err := m.DB.Prepare(sqlText)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		m.Entry.Name, // 姓名

		m.Entry.Age, //
	)
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
func (m *StudentModel) Find(ctx context.Context, id ...int64) (result []*Student, err error) {
	columns := strings.Join(m.columns, ",")
	if columns == "" {
		columns = m.getColumns()
	}

	whereField := m.getWhere()
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
func (m *StudentModel) Set(fields string, val interface{}) *StudentModel {
	m.set = append(m.set, tools.AddQuote(fields)+"=?")
	m.val = append(m.val, val)
	return m
}

// 更新数据
func (m *StudentModel) Update(ctx context.Context) (b bool, err error) {
	set := strings.Join(m.set, ",")
	whereField := m.getWhere()

	sqlText := "UPDATE " + " user " + " SET " + set + whereField

	return m.Save(ctx, sqlText, m.val...)
}

func (m *StudentModel) getWhere() string {
	if len(m.where) == 0 {
		return ""
	}
	return " where " + strings.Join(m.where, " and ")
}

func (m *StudentModel) getOr() string {
	if len(m.or) == 0 {
		return ""
	}
	return " or " + strings.Join(m.or, " or ")
}
