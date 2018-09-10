drop table if exists conversations_msgs;
drop table if exists conversations;
drop table if exists msgs;
drop table if exists users;

create table users
(
  id        serial primary key,
  tg_userid integer unique not null,
  tg_chatid integer        not null,
  type      varchar,
  name      text
);

create table msgs
(
  id       serial primary key                                                                        not null,
  userid   integer references users (id)                                                             not null,
  ts       timestamp,
  type     text                                                                                      not null,
  text     text,
  attachId bytea,
  caption  text
);

create table conversations
(
  id             integer primary key,
  description    text,
  status         int,
  creator        integer not null,
  created        timestamp,
  lastQuestionTs timestamp,
  lastAnswerTs   timestamp,
  totalMsgs      int,
  sla            timestamp,
  closedBy       int
);

create table conversations_msgs
(
  id  serial references conversations (id),
  msg int references msgs (id)
);

-- insert into users (tg_userid, tg_chatid, type, name) values (624938764, 624938764, 'Agent', 'pixel2');
-- insert into users (tg_userid, tg_chatid, type, name) values (55281426, 55281426, 'Agent', 'abefimov');
