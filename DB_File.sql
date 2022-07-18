drop database db;
create database db;
use db;

select * from users;

select * from organizations;

select * from jobs;

select * from categories;

select * from required_skills;

select * from profiles;

select * from skills;

select * from experiences;

show databases;

create table users(
id varchar (40) primary key,
email varchar (30),
name varchar (20),
pass varchar (20)
)


SELECT email, name, pass FROM user_data where id = "23ea5f3f-165d-43ff-960a-7e459ab301dd";

create table organizations(
org_id varchar (40) primary key,
name varchar (40),
about varchar (200),
website varchar (30),
u_id varchar(40),
FOREIGN KEY (u_id) REFERENCES users(id)
)


select * from organizations;

create table membership(
pk varchar (40),
id varchar(40),
org_id varchar(40),
FOREIGN KEY (id) REFERENCES users(id),
FOREIGN KEY (org_id) REFERENCES organizations(org_id)
)

select * from memberships;

insert into memberships values("a", "b", "c");




