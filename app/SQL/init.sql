DROP DATABASE IF EXISTS yourdbname; -- TODO remove this line in prod

CREATE DATABASE yourdbname;
USE yourdbname;

-- Create table 'users'
CREATE TABLE users (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    AESkey VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL
);

-- Create table 'stickies'
CREATE TABLE stickies (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    sticky_description VARCHAR(500) NOT NULL,
    sticky_title VARCHAR(255) NOT NULL,
    to_delete BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create table 'cards'
CREATE TABLE cards (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    card_id INT NOT NULL,
    balance VARCHAR(255) NOT NULL,
    due_date VARCHAR(255) NOT NULL,
    to_delete BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create table 'list_banks'
CREATE TABLE list_banks (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    bank_name VARCHAR(255) NOT NULL,
    bank_artwork VARCHAR(255) NOT NULL,
    to_delete BOOLEAN NOT NULL
);

-- Create table 'list_bank_accounts'
CREATE TABLE list_bank_accounts (
	id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	bank_id INT NOT NULL,
	amount VARCHAR(255) NOT NULL,
	to_delete BOOLEAN NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (bank_id) REFERENCES list_banks(id)
);

-- Create table 'CD'
CREATE TABLE CD (
	id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	bank_id INT NOT NULL,
	start_date VARCHAR(255) NOT NULL,
	deposit VARCHAR(255) NOT NULL,
	term VARCHAR(255) NOT NULL,
	apy VARCHAR(255) NOT NULL,
	to_delete BOOLEAN NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (bank_id) REFERENCES list_banks(id)
);

-- Create table 'list_cards'
CREATE TABLE list_cards (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    bank_id INT NOT NULL,
    card_name VARCHAR(255) NOT NULL,
    card_artwork VARCHAR(255) NOT NULL,
    FOREIGN KEY (bank_id) REFERENCES list_banks(id)
);
