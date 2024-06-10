USE yourdbname;

-- Drop tables if they exist
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS stickies;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS list_cards;
DROP TABLE IF EXISTS list_banks;

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
    bank_artwork VARCHAR(255) NOT NULL
);

-- Create table 'list_cards'
CREATE TABLE list_cards (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    bank_id INT NOT NULL,
    card_name VARCHAR(255) NOT NULL,
    card_artwork VARCHAR(255) NOT NULL,
    FOREIGN KEY (bank_id) REFERENCES list_banks(id)
);
