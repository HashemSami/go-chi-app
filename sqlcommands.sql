drop table if exists users

-- creating users table
create table users(
  id SERIAL PRIMARY KEY,
  age INT,
  first_name TEXT,
  last_name TEXT,
  email TEXT UNIQUE NOT NULL
);

-- we can skip the id because we set it as a serial key
insert into users (age, email, first_name, last_name)
values(30, 'hash@hash.com', 'hashem', 'sami')