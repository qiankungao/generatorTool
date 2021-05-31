package model

import (
	"database/sql"
	"fmt"
	"strconv"
)

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

