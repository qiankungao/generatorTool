CREATE TABLE `food`
(
    `id`         int(10) unsigned        NOT NULL AUTO_INCREMENT COMMENT '自增ID，商品Id',
    `name`       varchar(30)             NOT NULL COMMENT '商品名',
    `price`      decimal(10, 2) unsigned NOT NULL COMMENT '商品价格',
    `type_id`    int(10) unsigned        NOT NULL COMMENT '商品类型Id',
    `createtime` int(10)                 NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8