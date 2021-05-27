package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type UserModel struct {
	Entry     *User
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

func NewUser(db *sql.DB) *UserModel {
	return &UserModel{
		DB:        db,
		Entry:     &User{},
		columnMap: make(map[string]interface{}, 0),
	}
}

func NewUserWithEntry(db *sql.DB, entry *User) *UserModel {
	return &UserModel{
		DB:        db,
		Entry:     entry,
		columnMap: make(map[string]interface{}, 0),
	}
}

func NewUserTx(tx *sql.Tx) *UserModel {
	return &UserModel{
		Tx:    tx,
		Entry: &User{},
	}
}

// 获取所有的表字段
func (m *UserModel) getColumns() string {
	return " `id`,`username`,`age` "
}

//设置查询字段
func (m *UserModel) Columns(column ...string) *UserModel {
	for _, c := range column {
		m.columns = append(m.columns, c)

		if c == "id" {
			m.columnMap[c] = &m.Entry.Id
		}

		if c == "userName" {
			m.columnMap[c] = &m.Entry.UserName
		}

		if c == "age" {
			m.columnMap[c] = &m.Entry.Age
		}

	}
	return m
}

//设置查询字段
func (m *UserModel) Where(field string, val interface{}) *UserModel {
	m.where = append(m.where, field)
	m.val = append(m.val, val)
	return m
}

//设置查询字段
func (m *UserModel) And(field string, val interface{}) *UserModel {
	m.where = append(m.where, field)
	m.val = append(m.val, val)
	return m
}

//设置Or
func (m *UserModel) Or(field string, val interface{}) *UserModel {
	m.or = append(m.or, field)
	m.val = append(m.val, val)
	return m
}

func (m *UserModel) SetPageSize(size int) *UserModel {
	m.pageSize = size
	return m
}
func (m *UserModel) Page(page int) *UserModel {
	//(page-1)*pageSize, pageSize
	offset := (page - 1) * m.pageSize
	m.limit = " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(m.pageSize)
	return m
}

//设置联合查询的条件
func (m *UserModel) Joins(join string) *UserModel {
	m.joins = " " + join + " "
	return m
}

//生成set、get方法

func (m *UserModel) SetId(id int32) *UserModel {
	m.Entry.Id = id
	return m
}
func (m *UserModel) GetId() int32 {
	return m.Entry.Id
}

func (m *UserModel) SetUserName(userName string) *UserModel {
	m.Entry.UserName = userName
	return m
}
func (m *UserModel) GetUserName() string {
	return m.Entry.UserName
}

func (m *UserModel) SetAge(age int32) *UserModel {
	m.Entry.Age = age
	return m
}
func (m *UserModel) GetAge() int32 {
	return m.Entry.Age
}

// 获取多行数据.
func (m *UserModel) getRows(db *sql.DB, sqlTxt string, params ...interface{}) (rowsResult []*User, err error) {
	query, err := db.Query(sqlTxt, params...)
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
		row := User{}
		err = query.Scan(dest...)
		if err != nil {
			fmt.Println(err)
			return
		}

		row.Id = m.Entry.Id // 主键递增

		row.UserName = m.Entry.UserName // 姓名

		row.Age = m.Entry.Age //

		rowsResult = append(rowsResult, &row)
	}
	return
}

//获取单行
func (m *UserModel) getRow(db *sql.DB, sqlText string, params ...interface{}) (rowResult *User, err error) {
	query := m.DB.QueryRow(sqlText, params...)
	row := User{}
	err = query.Scan(
		&row.Id, // 主键递增

		&row.UserName, // 姓名

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
func (m *UserModel) Save(sqlTxt string, value ...interface{}) (b bool, err error) {
	stmt, err := m.DB.Prepare(sqlTxt)
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
func (m *UserModel) SaveTx(tx *sql.Tx, sqlTxt string, value ...interface{}) (b bool, err error) {
	stmt, err := tx.Prepare(sqlTxt)
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
func (m *UserModel) Create() (lastId int64, err error) {
	const sqlText = "INSERT INTO " + "user" + " (`username`,`age`) VALUES (?,?)"
	stmt, err := m.DB.Prepare(sqlText)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		m.Entry.UserName, // 姓名

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
func (m *UserModel) Find(id ...int64) (result []*User, err error) {
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

	result, err = m.getRows(m.DB, sqlText, m.val...)
	return
}

//设置更新字段
func (m *UserModel) Set(fields string, val interface{}) *UserModel {
	m.set = append(m.set, AddQuote(fields)+"=?")
	m.val = append(m.val, val)
	return m
}

// 更新数据
func (m *UserModel) Update() (b bool, err error) {
	set := strings.Join(m.set, ",")
	whereField := m.getWhere()

	sqlText := "UPDATE " + " user " + " SET " + set + whereField

	return m.Save(sqlText, m.val...)
}

func (m *UserModel) getWhere() string {
	if len(m.where) == 0 {
		return ""
	}
	return " where " + strings.Join(m.where, " and ")
}

func (m *UserModel) getOr() string {
	if len(m.or) == 0 {
		return ""
	}
	return " or " + strings.Join(m.or, " or ")
}
