-- connect db;
create table if not exists user (
  user_id integer primary key auto_increment,
  username char(255) not null,
  email char(255) not null,
  pw_hash char(255) not null
);

create table if not exists follower(
  who_id integer,
  whom_id integer
);

create table if not exists message (
  message_id integer primary key auto_increment,
  author_id integer not null,
  text varchar(5000) not null,
  pub_date integer,
  flagged integer
);
