USE yourdbname;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS stickies; 



CREATE TABLE users (
        id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
        username VARCHAR(5000) NOT NULL,
        password VARCHAR(5000) NOT NULL,
	email VARCHAR(5000) NOT NULL
);


CREATE TABLE stickies (
	id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	sticky_description VARCHAR(5000) NOT NULL,
	sticky_title VARCHAR(5000) NOT NULL,
	salt VARCHAR(5000) NOT NULL,
	to_delete BOOLEAN NOT NULL
);
