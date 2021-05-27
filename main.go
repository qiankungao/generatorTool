package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

/*

	func (s *{{.Name}}Select) Query(ctx context.Context, db *sql.DB) ([]*game.{{.Name}}, error) {
	sqlStr, args, err := s.handler.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	results := make([]*game.{{.Name}}, 0)
	columns, _ := rows.Columns()
	dest := make([]interface{}, 0, len(columns))

	for _, v := range columns {
		dest = append(dest, s.fieldMap[v])
	}
	for rows.Next() {
		result := &game.{{.Name}}{}
		err := rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		{{range $k, $v := .Cols}}
		result.{{$v.Name}} = s.tmp.{{$v.Name}}
		{{end}}
		results = append(results, result)
	}
	return results, nil
}


*/

func main() {

	//GenCurd()
	Run()
}

func Run() error {
	var err error

	DbConn = DBConfig{
		Host:        "127.0.0.1",
		Port:        3306,
		Name:        "root",
		Pass:        "root",
		DBName:      "helix2_game",
		Charset:     "utf8",
		MaxIdleConn: 5,
		MaxOpenConn: 10,
	}
	db, err := InitDB(DbConn)
	if db == nil || err != nil {
		return errors.New("database connect failed>>" + err.Error())
	}
	update(db)
	return nil
}

func update(db *sql.DB) {
	res,_:=NewUser(db).Set("userName","qiankun").Set("age",30).Where("Id = ?",1).Update()
	fmt.Println("res:",res)
}

func column(db *sql.DB) {
	user, _ := NewUser(db).Columns("userName", "age").Where("userName = ?", "gao").Find()

	for _, u := range user {
		fmt.Println("user:", u.UserName, u.Age)
	}
}

//or
func Or(db *sql.DB) {
	user, _ := NewUser(db).Where("userName = ?", "gao").Or("age = ?", 28).Find()

	for _, u := range user {
		fmt.Println("user:", u.Id)
	}
}

//联合查询
func Join(db *sql.DB) {
	user, _ := NewUser(db).Columns("user.userName", "user.age").Joins("left join email on user.Id=email.userId").Find()

	for _, u := range user {
		fmt.Println("user:", u.Id)
	}

}

//分页
func GetPage(db *sql.DB) {
	user, _ := NewUser(db).SetPageSize(10).Page(2).Find()

	for _, u := range user {
		fmt.Println("user:", u.Id)
	}

}

//查询
func GetRows(db *sql.DB) {
	user, _ := NewUser(db).
		Where("userName = ?", "xiaoming").And("age = ?", 28).
		Find()

	fmt.Println("user:", user)
}

//插入
func insert(db *sql.DB) {
	for i := 10; i < 50; i++ {
		ii := strconv.Itoa(i)
		id, err := NewUser(db).SetUserName("gao" + ii).SetAge(int32(i)).Create()
		fmt.Println("插入的元素：", id, err)
	}

}

//func command() {
//	//初使工作
//	DbModel := NewDB()
//	DbModel.Using(Conn)
//	DbModel.DBName = DbConn.DBName
//
//	dir, _ := os.Getwd()
//	logic := &Logic{
//		DB:   DbModel,
//		Path: dir + "/" + DefaultSavePath + DS, //默认当前命令所在目录
//	}
//	err := logic.DB.GetTableNameAndComment()
//	if err != nil {
//
//	}
//
//	commands := NewCommands(logic)
//	formats = []string{"json", "db"}
//	commands.GenerateEntry([]string{})
//	commands.GenerateCURD([]string{})
//	//commands.MarkDown([]string{})
//}
