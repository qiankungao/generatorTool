package entry

type SqlData struct {
	DescList []*SqlDataChild
}

type SqlDataChild struct {
	Index      int    // 自增
	TableName  string // 表名
	Comment    string // 表备注
	PrimaryKey string //主键
	List       []*SqlDesc
}

// 表结构详情
type SqlDesc struct {
	Index   int
	Name    string // 字段名字
	Type    string //字段类型
	Tag     string //tag
	Comment string // 备注
}

type CurdSqlInfo struct {
	TableName           string              // 表名
	PrimaryKey          string              // 主键字段
	PrimaryType         string              // 主键类型
	StructTableName     string              // 结构表名称
	NullStructTableName string              // 判断为空的表名
	PkgEntity           string              // 实体空间名称
	PkgTable            string              // 表的空间名称
	UpperTableName      string              // 大写的表名
	AllFieldList        string              // 所有字段列表,如: id,name
	InsertFieldList     string              // 插入字段列表,如:id,name
	InsertMark          string              // 插入字段使用多少个?,如: ?,?,?
	UpdateFieldList     string              // 更新字段列表
	SecondField         string              // 存放第二个字段
	UpdateListField     []string            // 更新字段列表
	FieldsInfo          []*SqlFieldInfo     // 字段信息
	NullFieldsInfo      []*NullSqlFieldInfo // 判断为空时
	InsertInfo          []*SqlFieldInfo
}

// 查询使用的字段结构信息
type SqlFieldInfo struct {
	HumpName  string // 驼峰字段名称
	Name      string //首字母小写字段
	FieldType string //字段的类型
	NullType  string
	Comment   string // 字段注释
}

type NullSqlFieldInfo struct {
	GoType       string // golang类型
	HumpName     string // 驼峰字段名称
	OriFieldType string // 原数据库类型
	Comment      string // 字段注释
}

var GoTypeToMysqlType = map[string]string{
	"string":  "varchar(255)",
	"int8":    "tinyint",
	"int16":   "smallint",
	"int32":   "integer",
	"int64":   "bigint",
	"int":     "int",
	"uint8":   "tinyint unsigned",
	"uint16":  "smallint unsigned",
	"uint32":  "integer unsigned",
	"uint64":  "bigint unsigned",
	"float32": "float",
	"float64": "double",
}

//MYSQL => golang mysql NULL TYPE
var GoTypeToMysqlNullType = map[string]string{
	"int8":    "sql.NullInt32",
	"int16":   "sql.NullInt32",
	"int":     "sql.NullInt32",
	"int32":   "sql.NullInt32",
	"int64":   "sql.NullInt64",
	"float32": "sql.NullFloat64",
	"float64": "sql.NullFloat64",
	"string":  "sql.NullString",
}

//mysql类型 <=> golang类型
var MysqlTypeToGoType = map[string]string{
	"tinyint":    "int32",
	"smallint":   "int32",
	"mediumint":  "int32",
	"int":        "int32",
	"integer":    "int64",
	"bigint":     "int64",
	"float":      "float64",
	"double":     "float64",
	"decimal":    "float64",
	"date":       "string",
	"time":       "string",
	"year":       "string",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"char":       "string",
	"varchar":    "string",
	"tinyblob":   "string",
	"tinytext":   "string",
	"blob":       "string",
	"text":       "string",
	"mediumblob": "string",
	"mediumtext": "string",
	"longblob":   "string",
	"longtext":   "string",
}

//MYSQL => golang mysql NULL TYPE
var MysqlTypeToGoNullType = map[string]string{
	"tinyint":    "sql.NullInt32",
	"smallint":   "sql.NullInt32",
	"mediumint":  "sql.NullInt32",
	"int":        "sql.NullInt32",
	"integer":    "sql.NullInt64",
	"bigint":     "sql.NullInt64",
	"float":      "sql.NullFloat64",
	"double":     "sql.NullFloat64",
	"decimal":    "sql.NullFloat64",
	"date":       "sql.NullString",
	"time":       "sql.NullString",
	"year":       "sql.NullString",
	"datetime":   "mysql.NullTime",
	"timestamp":  "mysql.NullTime",
	"char":       "sql.NullString",
	"varchar":    "sql.NullString",
	"tinyblob":   "sql.NullString",
	"tinytext":   "sql.NullString",
	"blob":       "sql.NullString",
	"text":       "sql.NullString",
	"mediumblob": "sql.NullString",
	"mediumtext": "sql.NullString",
	"longblob":   "sql.NullString",
	"longtext":   "sql.NullString",
}
