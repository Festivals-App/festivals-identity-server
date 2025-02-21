--
-- Create the Festivals Identity Database
--

-- First create the database
CREATE DATABASE IF NOT EXISTS `festivals_identity_database`;

-- Create the tables in the newly created database
USE festivals_identity_database;

/**
Create the basic entities
*/

-- Create the users table
CREATE TABLE IF NOT EXISTS `users` (

	`user_id` 			    int unsigned 	 	NOT NULL AUTO_INCREMENT 											    COMMENT 'The id of the user.',
	`user_email` 		    varchar(255)		NOT NULL													            COMMENT 'The email of the user. The email needs to be unique.',
	`user_password` 	    varchar(225) 	  	NOT NULL 												                COMMENT 'The password hash of the users password.',
	`user_createdat` 		timestamp 			NOT NULL DEFAULT current_timestamp()					      		    COMMENT 'The date and time the user was created.',
	`user_updatedat` 		timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()	    COMMENT 'The date and time the user data was last updated.',
    `user_role` 	  	    tinyint 		    NOT NULL DEFAULT 0											            COMMENT 'The role of the user.',

PRIMARY 	KEY (`user_id`),
UNIQUE 	    KEY (`user_email`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='The user table represents a user that interacts with the FestivalsApp backend.';

-- Create the service key table
CREATE TABLE IF NOT EXISTS `service_keys` (

	`service_key_id` 	    int unsigned 	 	NOT NULL AUTO_INCREMENT 									    COMMENT 'The id of the key.',
	`service_key` 	  	    varchar(225) 		NOT NULL 												        COMMENT 'The service key.',
    `service_key_comment` 	varchar(225) 		NOT NULL 												        COMMENT 'A comment about the service key.',

PRIMARY 	KEY (`service_key_id`),
UNIQUE 	  	KEY (`service_key`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table contains all service node keys.';

-- Create the service key table
CREATE TABLE IF NOT EXISTS `api_keys` (

	`api_key_id` 			int unsigned 	 	NOT NULL AUTO_INCREMENT 											COMMENT 'The id of the key.',
	`api_key` 	  	        varchar(225) 		NOT NULL 												            COMMENT 'The api key.',
    `api_key_comment` 	  	varchar(225) 		NOT NULL 												            COMMENT 'A comment about the api key.',

PRIMARY 	KEY (`api_key_id`),
UNIQUE 	  	KEY (`api_key`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table contains all api keys.';

/**
Create the mapping tables to associate entities
*/

-- Create the table to map festivals to users
CREATE TABLE IF NOT EXISTS `map_festival_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_festival` 		int unsigned 		NOT NULL					        COMMENT 'The id of the mapped festival.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_festival`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps festivals to users.';

-- Create the table to map artists to users
CREATE TABLE IF NOT EXISTS `map_artist_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_artist` 		int unsigned 		NOT NULL					        COMMENT 'The id of the mapped artist.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_artist`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps artists to users.';

-- Create the table to map locations to users
CREATE TABLE IF NOT EXISTS `map_location_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_location` 		int unsigned 		NOT NULL					        COMMENT 'The id of the mapped location.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_location`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps locations to users.';

-- Create the table to map events to users
CREATE TABLE IF NOT EXISTS `map_event_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_event` 		    int unsigned 		NOT NULL					        COMMENT 'The id of the mapped event.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_event`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps events to users.';

-- Create the table to map links to users
CREATE TABLE IF NOT EXISTS `map_link_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_link` 		    int unsigned 		NOT NULL					        COMMENT 'The id of the mapped link.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_link`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps links to users.';

-- Create the table to map images to users
CREATE TABLE IF NOT EXISTS `map_image_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_image` 		    int unsigned 		NOT NULL					        COMMENT 'The id of the mapped image.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_image`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps images to users.';

-- Create the table to map places to users
CREATE TABLE IF NOT EXISTS `map_place_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_place` 		    int unsigned 		NOT NULL					        COMMENT 'The id of the mapped place.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_place`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps places to users.';

-- Create the table to map tags to users
CREATE TABLE IF NOT EXISTS `map_tag_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_tag` 		    int unsigned 		NOT NULL					        COMMENT 'The id of the mapped tag.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_tag`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps tags to users.';

/**
Insert default users (default password: we4711), api key and service key.
*/

INSERT INTO  `users`(`user_id`, `user_email`, `user_password`, `user_role`) VALUES (0, 'admin@email.com', '$2a$12$YbAhewILx82tGkLtEZWiKOfYzBt85RSQtGXhxlQX2hV7qiP51xPES', 42);
INSERT INTO  `users`(`user_id`, `user_email`, `user_password`, `user_role`) VALUES (0, 'user@email.com', '$2a$12$YbAhewILx82tGkLtEZWiKOfYzBt85RSQtGXhxlQX2hV7qiP51xPES', 1);
INSERT INTO  `users`(`user_id`, `user_email`, `user_password`, `user_role`) VALUES (0, 'coordinator@email.com', '$2a$12$YbAhewILx82tGkLtEZWiKOfYzBt85RSQtGXhxlQX2hV7qiP51xPES',2);
INSERT INTO api_keys(`api_key`, `api_key_comment`) VALUES ('TEST_API_KEY_001', "DEVELOPMENT API KEY");
INSERT INTO service_keys(`service_key`, `service_key_comment`) VALUES ('TEST_SERVICE_KEY_001', "DEVELOPMENT SERVICE KEY");