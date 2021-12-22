-- 如果已经存在这个数据库则删除它
drop database if exists `ManageSystem`; 
-- 重新创建数据库
create database `ManageSystem`;
-- 创建用户信息表(user)
use `ManageSystem`;

create table `user`(
    `user_id` int not null AUTO_INCREMENT,  -- not null默认不为空 AUTO_INCREMENT自增定义
    `user_name` varchar(100) not null Unique, -- varchar用于存储可变长字符串,并将user_name作为唯一索引，来实现不重复插入
    `user_password` varchar(100) not null,
    `user_role` int not null, -- 用于存储用户的等级
    `user_age` int,
    `user_gender` varchar(100),
    PRIMARY KEY (user_id)
);

create table `player`(
    `player_id` int not null AUTO_INCREMENT ,
    `player_name` varchar(100) not null unique,
    `player_team` varchar(100),
    `player_avatar` varchar(100),
    `player_information` varchar(100),
    PRIMARY KEY(player_id)
);
create table `team`(
    `team_id` int not null AUTO_INCREMENT,
    `team_name` varchar(100),
    `team_logo` varchar(100),
    `team_member` int,
    PRIMARY KEY(team_id)
);


create table `match`( -- 球赛表
    `match_id` int not null AUTO_INCREMENT,
	`match_name` varchar(100) not null Unique,
	`match_date` date not null,
	`match_place` varchar(100) not null,
	`match_info` varchar(100)  not null comment '详情',
	`match_appointment` int comment '预约数',
	`match_teama` varchar(100),
	`match_teamb` varchar(100),
    PRIMARY KEY(match_id)
) comment='球赛服务';

create table `usertomatch`( -- 用户和球赛的关系表
    `usertomatch_id` int not null AUTO_INCREMENT,
    `user_name` varchar(100) not null,
    `match_name` varchar(100) not null,
    PRIMARY KEY(usertomatch_id)
);



create table `playertoteam`(
    `player_name` int not NULL,
    `team_name` varchar(100),
    PRIMARY KEY(player_id)
);