 {{range $key, $item := .DescList}}
DROP TABLE IF EXISTS `{{$item.TableName}}`;

 CREATE TABLE `{{$item.TableName}}` (
 {{range $item.List}} `{{.Name}}` {{.Type}} {{.Tag}} {{.Comment}},
 {{end}}
  PRIMARY KEY (`{{$item.PrimaryKey}}`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT {{$item.Comment}};
 {{end}}