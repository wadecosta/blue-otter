CREATE TABLE events (
	id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	event_data VARCHAR(10000) NOT NULL,
	event_when VARCHAR(5000) NOT NULL
);