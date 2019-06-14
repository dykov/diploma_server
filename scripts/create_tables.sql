create table users (
  id bigserial primary key ,
  first_name varchar(100) not null ,
  middle_name varchar(100) not null ,
  last_name varchar(100) not null ,
  login varchar(100) unique not null ,
  password varchar(150) not null ,
  email varchar(300) unique not null ,
  occupation varchar(200) ,
  is_onaft_student int default 0 ,
  rating bigint default 0 ,
  role int default 0 ,
  verification_code varchar(30)
);

create table courses (
  id bigserial primary key ,
  name text not null
);

create table sections (
  id bigserial primary key ,
  course_id bigint references courses(id) ,
  name text not null
);

create table lessons (
  id bigserial primary key ,
  section_id bigint references sections(id) ,
  name text not null
);

create table paragraphs_or_tests (
  id bigserial primary key ,
  lesson_id bigint references lessons(id) ,
  name text not null ,
  text text not null ,
  points int default 0
);

create table tests_answers (
  id bigserial primary key ,
  test_id bigint references paragraphs_or_tests(id) ,
  text text not null ,
  is_right_answer boolean default false
);

create table users_tests (
  user_id bigint references users(id),
  test_id bigint references paragraphs_or_tests(id)
);

create table backups (
  time timestamp default now(),
  filepath text,
  name text
);