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

-- Create the festival table
CREATE TABLE IF NOT EXISTS `users` (

	`user_id` 			    int unsigned 	 	NOT NULL AUTO_INCREMENT 											    COMMENT 'The id of the user.',
	`user_name` 	  	    varchar(225) 		NOT NULL 												                COMMENT 'The name of the user. The name needs to be unique.',
	`user_email` 		    varchar(255)		NOT NULL													            COMMENT 'The email of the user. The email needs to be unique.',
	`user_password` 	    varchar(225) 	  	NOT NULL 												                COMMENT '',
	`user_tokenhash` 		varchar(15) 	  	NOT NULL 											                    COMMENT '',
	`user_createdat` 		timestamp 			NOT NULL DEFAULT current_timestamp()					      		    COMMENT '',
	`user_updatedat` 		timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()	    COMMENT '',
    `user_role` 	  	    tinyint 		    NOT NULL DEFAULT 0											            COMMENT 'The role of the user.',

PRIMARY 	KEY (`user_id`),
UNIQUE 	  	KEY `name` (`user_name`),
            KEY `email` (`user_email`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='The user table represents a user that interacts with the FestivalsApp backend.';

/**

Create the mapping tables to associate entities

*/

-- Create the table to map festivals to users
CREATE TABLE IF NOT EXISTS `map_festival_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_festival` 		int unsigned 		NOT NULL					        COMMENT 'The id of the mapped festival.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps festivals to users.';

-- Create the table to map artists to users
CREATE TABLE IF NOT EXISTS `map_artist_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_artist` 		int unsigned 		NOT NULL					        COMMENT 'The id of the mapped artist.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps artists to users.';

-- Create the table to map locations to users
CREATE TABLE IF NOT EXISTS `map_location_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_location` 		int unsigned 		NOT NULL					        COMMENT 'The id of the mapped location.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps locations to users.';