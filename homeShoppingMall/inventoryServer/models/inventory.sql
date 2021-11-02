DROP TABLE if EXISTS `inventory`
CREATE TABLE `inventory`
(
    `id`            bigint(32) primary key auto_increment,
    `goods`         varchar(32) not null,
    `inventory_num` bigint(64) not null,
    UNIQUE KEY `inx_goods` (`goods`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE if EXISTS `inventory_shell`;
CREATE TABLE `inventory_shell`
(
    `id`         bigint(32) primary key auto_increment,
    `orders_goods_num` VARCHAR(64),
    `status`     ENUM('NO', 'OK', 'FAILED') default 'NO',
    UNIQUE KEY `inx_orders_num` (`orders_goods_num`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE if EXISTS `inventory_shell_detail`;
CREATE TABLE `inventory_shell_detail`
(
    `id`    bigint(32) primary key auto_increment,
    `goods` VARCHAR(64),
    `num`   INT(32),
    `orders_goods_num` VARCHAR(64),
    foreign key (orders_goods_num) references `inventory_shell` (orders_goods_num) on update cascade on delete cascade
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;