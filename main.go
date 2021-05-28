package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/1975210542/generatorTools/db"
	"github.com/1975210542/generatorTools/generator"
)

var (
	scanPath      string
	outputPath    string
	generatorType string
)

func init() {
	flag.StringVar(&scanPath, "scanPath", "model", "scanPath")
	flag.StringVar(&outputPath, "outputPath", "/output/mysql", "outputPath")
	flag.StringVar(&generatorType, "generatorType", "db", "generatorType")
}

func main() {
	flag.Parse()
	fmt.Println("scanPath:", scanPath)
	fmt.Println("scanPath:", outputPath)

	switch generatorType {
	case "sql":
		generator.GeneratorSql(scanPath, outputPath)
	case "db":
		generator.GeneratorCurd(scanPath, outputPath)
	}

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
