//
package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/1975210542/generatorTools/output/db_models/config"
	"github.com/1975210542/generatorTools/output/entity"
	_ "github.com/go-sql-driver/mysql"
)

type BorderModel struct {
	DB *sql.DB
	Tx *sql.Tx
}

// not transaction
func NewBorder(db ...*sql.DB) *BorderModel {
	if len(db) > 0 {
		return &BorderModel{
			DB: db[0],
		}
	}
	return &BorderModel{
		//DB: masterDB,
	}
}

// transaction object
func NewBorderTx(tx *sql.Tx) *BorderModel {
	return &BorderModel{
		Tx: tx,
	}
}

// 获取所有的表字段
func (m *BorderModel) getColumns() string {
	return " `id`,`playerId`,`stageInfo`,`freeTimes` "
}

// 获取多行数据.
func (m *BorderModel) getRows(sqlTxt string, params ...interface{}) (rowsResult []*entity.Border, err error) {
	query, err := m.DB.Query(sqlTxt, params...)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer query.Close()
	for query.Next() {
		row := entity.BorderNull{}
		err = query.Scan(
			&row.Id,        //
			&row.PlayerId,  //
			&row.StageInfo, // status:关卡的状态0 可挑战 1 可放置 2 已放置； model:模式 1 简单 2 困难；hangRoles:驻扎的英雄； endTime:驻扎的结束时间
			&row.FreeTimes, // 免费的重置次数
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		rowsResult = append(rowsResult, &entity.Border{
			Id:        row.Id.Int32,         //
			PlayerId:  row.PlayerId.Int32,   //
			StageInfo: row.StageInfo.String, // status:关卡的状态0 可挑战 1 可放置 2 已放置； model:模式 1 简单 2 困难；hangRoles:驻扎的英雄； endTime:驻扎的结束时间
			FreeTimes: row.FreeTimes.Int32,  // 免费的重置次数
		})
	}
	return
}

// 获取单行数据
func (m *BorderModel) getRow(sqlText string, params ...interface{}) (rowResult *entity.Border, err error) {
	query := m.DB.QueryRow(sqlText, params...)
	row := entity.BorderNull{}
	err = query.Scan(
		&row.Id,        //
		&row.PlayerId,  //
		&row.StageInfo, // status:关卡的状态0 可挑战 1 可放置 2 已放置； model:模式 1 简单 2 困难；hangRoles:驻扎的英雄； endTime:驻扎的结束时间
		&row.FreeTimes, // 免费的重置次数
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	rowResult = &entity.Border{
		Id:        row.Id.Int32,         //
		PlayerId:  row.PlayerId.Int32,   //
		StageInfo: row.StageInfo.String, // status:关卡的状态0 可挑战 1 可放置 2 已放置； model:模式 1 简单 2 困难；hangRoles:驻扎的英雄； endTime:驻扎的结束时间
		FreeTimes: row.FreeTimes.Int32,  // 免费的重置次数
	}
	return
}

// _更新数据
func (m *BorderModel) Save(sqlTxt string, value ...interface{}) (b bool, err error) {
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

// _更新数据
func (m *BorderModel) SaveTx(sqlTxt string, value ...interface{}) (b bool, err error) {
	stmt, err := m.Tx.Prepare(sqlTxt)
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

// 新增信息
func (m *BorderModel) Create(value *entity.Border) (lastId int64, err error) {
	const sqlText = "INSERT INTO " + config.TABLE_BORDER + " (`playerId`,`stageInfo`,`freeTimes`) VALUES (?,?,?)"
	stmt, err := m.DB.Prepare(sqlText)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		value.PlayerId,  //
		value.StageInfo, // status:关卡的状态0 可挑战 1 可放置 2 已放置； model:模式 1 简单 2 困难；hangRoles:驻扎的英雄； endTime:驻扎的结束时间
		value.FreeTimes, // 免费的重置次数
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

// 新增信息 tx
func (m *BorderModel) CreateTx(value *entity.Border) (lastId int64, err error) {
	const sqlText = "INSERT INTO " + config.TABLE_BORDER + " (`playerId`,`stageInfo`,`freeTimes`) VALUES (?,?,?)"
	stmt, err := m.Tx.Prepare(sqlText)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		value.PlayerId,  //
		value.StageInfo, // status:关卡的状态0 可挑战 1 可放置 2 已放置； model:模式 1 简单 2 困难；hangRoles:驻扎的英雄； endTime:驻扎的结束时间
		value.FreeTimes, // 免费的重置次数
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

// 更新数据
func (m *BorderModel) Update(value *entity.Border) (b bool, err error) {
	const sqlText = "UPDATE " + config.TABLE_BORDER + " SET `playerId`=?,`stageInfo`=?,`freeTimes`=? WHERE `id` = ?"
	params := make([]interface{}, 0)
	params = append(params, value.PlayerId)
	params = append(params, value.StageInfo)
	params = append(params, value.FreeTimes)
	params = append(params, value.Id)

	return m.Save(sqlText, params...)
}

// 更新数据 tx
func (m *BorderModel) UpdateTx(value *entity.Border) (b bool, err error) {
	const sqlText = "UPDATE " + config.TABLE_BORDER + " SET `playerId`=?,`stageInfo`=?,`freeTimes`=? WHERE `id` = ?"
	params := make([]interface{}, 0)
	params = append(params, value.PlayerId)
	params = append(params, value.StageInfo)
	params = append(params, value.FreeTimes)
	params = append(params, value.Id)

	return m.SaveTx(sqlText, params...)
}

// 查询多行数据
func (m *BorderModel) Find() (resList []*entity.Border, err error) {
	sqlText := "SELECT" + m.getColumns() + "FROM " + config.TABLE_BORDER
	resList, err = m.getRows(sqlText)
	return
}

// 获取单行数据
func (m *BorderModel) First(id int64) (result *entity.Border, err error) {
	sqlText := "SELECT" + m.getColumns() + "FROM " + config.TABLE_BORDER + " WHERE `id` = ? LIMIT 1"
	result, err = m.getRow(sqlText, id)
	return
}

// 获取最后一行数据
func (m *BorderModel) Last() (result *entity.Border, err error) {
	sqlText := "SELECT" + m.getColumns() + "FROM " + config.TABLE_BORDER + " ORDER BY ID DESC LIMIT 1"
	result, err = m.getRow(sqlText)
	return
}

// 单列数据
func (m *BorderModel) Pluck(id int64) (result map[int64]interface{}, err error) {
	const sqlText = "SELECT `id`, `playerId` FROM " + config.TABLE_BORDER + " where `id` = ?"
	rows, err := m.DB.Query(sqlText, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	result = make(map[int64]interface{})
	var (
		_id  int64
		_val interface{}
	)
	for rows.Next() {
		err = rows.Scan(&_id, &_val)
		if err != nil {
			fmt.Println(err)
			return
		}
		result[_id] = _val
	}
	return
}

// 单列数据 by 支持切片传入
// Get column data
func (m *BorderModel) Plucks(ids []int64) (result map[int64]interface{}, err error) {
	result = make(map[int64]interface{})
	if len(ids) == 0 {
		return
	}
	sqlText := "SELECT `id`, `playerId` FROM " + config.TABLE_BORDER + " where " +
		"`id` in (" + RepeatQuestionMark(len(ids)) + ")"
	params := make([]interface{}, len(ids))
	for idx, id := range ids {
		params[idx] = id
	}
	rows, err := m.DB.Query(sqlText, params...)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	var (
		_id  int64
		_val interface{}
	)
	for rows.Next() {
		err = rows.Scan(&_id, &_val)
		if err != nil {
			fmt.Println(err)
			return
		}
		result[_id] = _val
	}
	return
}

// 获取单个数据
// Get one data
func (m *BorderModel) One(id int64) (result int64, err error) {
	sqlText := "SELECT `id` FROM " + config.TABLE_BORDER + " where `id`=?"
	err = m.DB.QueryRow(sqlText, id).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}
	return
}

// 获取行数
// Get line count
func (m *BorderModel) Count() (count int64, err error) {
	sqlText := "SELECT COUNT(*) FROM " + config.TABLE_BORDER
	err = m.DB.QueryRow(sqlText).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}
	return
}

// 判断数据是否存在
// Check the data is have?
func (m *BorderModel) Has(id int64) (b bool, err error) {
	sqlText := "SELECT `id` FROM " + config.TABLE_BORDER + " where `id` = ?"
	var count int64
	err = m.DB.QueryRow(sqlText, id).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}
	return count > 0, nil
}

// repeat response to ?,?,?
func RepeatQuestionMark(count int) string {
	return strings.TrimRight(strings.Repeat("?,", count), ",")
}
