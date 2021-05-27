DROP TABLE IF EXISTS `user`;

CREATE TABLE `user`
(
    `Id`     int          not null AUTO_INCREMENT COMMENT '主键递增',
    `userName` varchar(200) not null default 'gaoqiankun' COMMENT '姓名',
    `age`  int(20)      not null COMMENT '',

    PRIMARY KEY (`Id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT '用户表';

DROP TABLE IF EXISTS `Student`;

CREATE TABLE `Student`
(
    `Id`     int          not null AUTO_INCREMENT COMMENT '主键递增',
    `myName` varchar(200) not null default 'gaoqiankun' COMMENT '姓名',
    `myAge`  int(20)      not null COMMENT '',

    PRIMARY KEY (`myName`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT '学生表';
 