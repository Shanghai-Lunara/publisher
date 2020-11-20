CREATE DATABASE publisher
  CHARACTER SET utf8
  COLLATE utf8_general_ci;

CREATE TABLE records (
    id BIGINT NOT NULL AUTO_INCREMENT,
    PRIMARY KEY(id),
    namespace VARCHAR(128) DEFAULT '' comment 'namespace项目命名空间',
    groupName VARCHAR(128) DEFAULT '' comment '项目分支渠道名称',
    runnerName VARCHAR(128) DEFAULT '' comment 'runner名称',
    stepInfo BLOB comment '步骤完整结束时完整信息',
    createdTM INT(11) NOT NULL
);


