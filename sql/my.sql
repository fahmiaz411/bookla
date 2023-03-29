CREATE TABLE users (
	id int auto_increment primary key,
    phone varchar(16) unique key,
    username varchar(20) unique key,
    fullname varchar(50),
    email varchar(50) unique key,
    gender enum('male', 'female'),
    birthdate date,
    pin varchar(200),
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp,
    deleted_at datetime
);

ALTER TABLE users ADD column pin varchar(200) AFTER birthdate;

SELECT * FROM users;