DROP TABLE IF EXISTS `order_shopping_car`;
CREATE TABLE `order_shopping_car`
(
    `id`          bigint(64) PRIMARY KEY auto_increment,
    `user_id`     bigint(64) NOT NULL,
    `goods`       VARCHAR(64) NOT NULL,
    `goods_nums`  bigint(64) NOT NULL default 0,
    `selected`    bool        NOT NULL DEFAULT FALSE,
    `create_time` timestamp null default current_timestamp (),
    INDEX(`user_id`),
    INDEX(`goods`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `order_info`;
CREATE TABLE `order_info`
(
    `id`               bigint(64) PRIMARY KEY auto_increment,
    `user_id`          bigint(64) not null,
    `order_all_price`  bigint(64),
    `pay_type`         ENUM('支付宝', '微信', '其他') default '支付宝',
    `status`           ENUM('WAIT_BUYER_PAY', 'TRADE_CLOSED', 'TRADE_SUCCESS', 'TRADE_FINISHED') default 'WAIT_BUYER_PAY',
    `goods_order_num`  VARCHAR(64) not null,
    `alipay_order_num` VARCHAR(64) default '',
    `address`          VARCHAR(64),
    `name`             VARCHAR(64),
    `phone`            VARCHAR(64),
    `message`          VARCHAR(128) default '',
    `pay_time`         timestamp,
    `create_time`      timestamp null default current_timestamp (),
    UNIQUE KEY inx_goods_order_num (`goods_order_num`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `order_goods_info`;
CREATE TABLE `order_goods_info`
(
    `id`              bigint(64) PRIMARY KEY auto_increment,
    `goods_sell_num`  int,
    `goods_order_num` VARCHAR(64) not null,
    `goods`           VARCHAR(64),
    `goods_price`     VARCHAR(64),
    `create_time`     timestamp null default current_timestamp (),
    INDEX(`goods_order_num`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;