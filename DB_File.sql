drop database db;
create database db;
use db;

select * from user_data;


create table user_data(
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
website varchar (30)
)

select * from organizations;

create table membership(
pk varchar (40),
id varchar(40),
org_id varchar(40),
FOREIGN KEY (id) REFERENCES user_data(id),
FOREIGN KEY (org_id) REFERENCES organizations(org_id)
)

select * from membership;






