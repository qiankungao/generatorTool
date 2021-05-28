package model

//@ Comment 用户表
//@ Name user
type User struct {
	Id       int32  `db:"not null;AUTO_INCREMENT"`                              //主键递增
	UserName string `db:"column:myName;size:200;not null;default 'gaoqiankun'"` //姓名
	Age      int32  `db:column:myAge;size:20;not null`
}
