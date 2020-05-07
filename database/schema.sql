create table settings
(
    id    int auto_increment
        primary key,
    `key` varchar(255) not null,
    value text         not null,
    constraint settings_key_uindex
        unique (`key`)
) charset = utf8;

create table tasks
(
    id          int auto_increment
        primary key,
    title       varchar(255)                         not null,
    object_type varchar(50)                          not null,
    object_id   int                                  not null,
    `interval`  int                                  not null comment 'Переодичность запуска в минутах',
    next_time   timestamp  default CURRENT_TIMESTAMP not null,
    enabled     tinyint(1) default 0                 not null,
    usernames   text                                 null comment 'Строка юзернеймов разделенных запятой, кому отправить уведомление
Если null, бдует отправленно в дефолтный канал телеграмма'
) charset = utf8;

create table tasks_graylog
(
    id             int auto_increment
        primary key,
    pattern        text          not null,
    aggregate_time varchar(50)   not null,
    min            int default 0 not null,
    max            int default 0 not null
) charset = utf8;

create table results_graylog
(
    id              int auto_increment
        primary key,
    task_graylog_id int                                 not null,
    status          int                                 not null,
    title           varchar(150)                        not null,
    message         varchar(150)                        not null,
    text            text                                not null,
    count           int                                 not null,
    graylog_link    text                                not null,
    created_at      timestamp default CURRENT_TIMESTAMP not null,
    constraint results_graylog_tasks_graylog_id_fk
        foreign key (task_graylog_id) references tasks_graylog (id)
) charset = utf8;

create table telegram_chats
(
    id         int auto_increment
        primary key,
    username   varchar(255)                        not null,
    chat_id    varchar(255)                        not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    constraint bot_chats_chat_id_uindex
        unique (chat_id),
    constraint bot_chats_username_uindex
        unique (username)
) charset = utf8;

