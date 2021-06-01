DROP TABLE IF EXISTS `student`;

CREATE TABLE `student`
(
    `Id`     integer      not null AUTO_INCREMENT COMMENT '主键递增',
    `myName` varchar(200) not null default 'gaoqiankun' COMMENT '姓名',
    `myAge`  integer(20) not null COMMENT '',

    PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';

DROP TABLE IF EXISTS `student1`;

CREATE TABLE `student1`
(
    `Id`     integer                           not null AUTO_INCREMENT COMMENT '主键递增',
    `myName` varchar(200) default 'gaoqiankun' not null COMMENT '姓名',
    `myAge`  integer(20) not null COMMENT '',

    PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user`
(
    `Id`     integer      not null AUTO_INCREMENT COMMENT '主键递增',
    `myName` varchar(200) not null default 'gaoqiankun' COMMENT '姓名',
    `myAge`  integer(20) not null COMMENT '',

    PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';
 