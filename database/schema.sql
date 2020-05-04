create table tasks
(
    id          int auto_increment
        primary key,
    title       varchar(255)         not null,
    is_enable   tinyint(1) default 1 not null,
    object_type int        default 1 not null,
    object_id   int                  not null,
    `interval`  int                  null,
    time        varchar(50)          null
)
    charset = utf8;

create table tasks_graylog
(
    id        int auto_increment
        primary key,
    pattern   text         not null,
    agg_time  varchar(255) not null,
    min_count int          null,
    max_count int          null
)
    charset = utf8;

