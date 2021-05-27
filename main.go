package main

import (
	"errors"
	"github.com/1975210542/generatorTools/db"
	"github.com/1975210542/generatorTools/generator"
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

	generator.GenCurd()
	//generator.GeneratorSql()
	//Run()
}

func Run() error {
	var err error

	DbConn := db.DBConfig{
		Host:        "127.0.0.1",
		Port:        3306,
		Name:        "root",
		Pass:        "root",
		DBName:      "helix2_game",
		Charset:     "utf8",
		MaxIdleConn: 5,
		MaxOpenConn: 10,
	}
	db, err := db.InitDB(DbConn)
	if db == nil || err != nil {
		return errors.New("database connect failed>>" + err.Error())
	}
	//update(db)
	return nil
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
