drop table if exists `classify_goods`;
create table `classify_goods`
(
    `id`   bigint(32) primary key auto_increment,
    `name` varchar(32) not null,
    `pid`  bigint(64) not null default 0,
    UNIQUE KEY `inx_name` (`name`)
)CHARSET=utf8 COLLATE=utf8_bin;

drop table if exists `goods`;
create table `goods`
(
    `id`                 int primary key auto_increment,
    `name`               varchar(32) not null,
    `goods_price`        varchar(32) not null,
    `promotion_price`    varchar(32) not null default 0,
    `classify_name`      varchar(32) not null,
    `goods_introduction` varchar(64) not null,
    `sales_volume`       int         not null default 0,
    `collect_num`        int         not null default 0,
    `is_show`            bool                 default false,
    `is_new`             bool                 default false,
    `is_freight_free`    bool                 default false,
    `is_hot`             bool                 default false,
    `create_time`        timestamp null default current_timestamp (),
    UNIQUE KEY `inx_name` (`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

drop table if exists `goods_image`;
create table `goods_image`
(
    `id`         bigint(32) primary key auto_increment,
    `image_url`  varchar(64) not null,
    `goods_name` varchar(64),
    foreign key (goods_name) references goods (name) on update cascade on delete cascade
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

drop table if exists `banner`;
create table `banner`
(
    `id`              int primary key auto_increment,
    `image_url`       varchar(64) not null,
    `image_goods_url` varchar(64) not null,
    `level`           int default 0,
    UNIQUE KEY `inx_level` (`level`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;