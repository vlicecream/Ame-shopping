DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`          bigint(20) primary key auto_increment not null,
    `user_id`     bigint(64) not null,
    `nick_name`   varchar(32),
    `mobile`      varchar(32) not null,
    `password`    varchar(64) not null,
    `gender`      enum('male', 'female', 'other') default 'other',
    `role`        enum('user', 'admin') default 'user',
    `create_time` timestamp null default current_timestamp (),
    `delete_time` timestamp null default current_timestamp () on update current_timestamp (),
    `is_delete`   boolean     default false ,
    UNIQUE KEY `inx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TRIGGER foo BEFORE INSERT ON user FOR EACH ROW
    IF NEW.nick_name IS NULL THEN
        SET NEW.nick_name := NEW.user_id ;
END IF;;