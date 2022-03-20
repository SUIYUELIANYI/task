DROP DATABASE IF EXISTS `student`;

CREATE DATABASE `student`;

USE `student`;

-- 爬取21级的姓名(带*号)
CREATE TABLE `users`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `student_id` VARCHAR(100) UNIQUE NOT NULL,-- 学号
    `name` VARCHAR(100) NOT NULL,
    `grade` VARCHAR(100) NOT NULL,
    PRIMARY KEY(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

