package model

import "time"

/*
		`user_id` 			    int unsigned 	 	NOT NULL AUTO_INCREMENT 											    COMMENT 'The id of the user.',
		`user_name` 	  	    varchar(225) 		NOT NULL 												                COMMENT 'The name of the user. The name needs to be unique.',
		`user_email` 		    varchar(255)		NOT NULL													            COMMENT 'The email of the user. The email needs to be unique.',
		`user_password` 	    varchar(225) 	  	NOT NULL 												                COMMENT '',
		`user_tokenhash` 		varchar(15) 	  	NOT NULL 											                    COMMENT '',
		`user_createdat` 		timestamp 			NOT NULL DEFAULT current_timestamp()					      		    COMMENT '',
		`user_updatedat` 		timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()	    COMMENT '',
	    `user_role` 	  	    tinyint 		    NOT NULL DEFAULT 0											            COMMENT 'The role of the user.',
*/

type User struct {
	ID         int         `json:"user_id" sql:"user_id"`
	Name       string      `json:"user_name" sql:"user_name"`
	Email      string      `json:"user_email" sql:"user_email"`
	Password   string      `json:"user_password" sql:"user_password"`
	TokenHash  string      `json:"user_tokenhash" sql:"user_tokenhash"`
	CreateDate time.Time   `json:"user_createdat" sql:"user_createdat"`
	UpdateDate time.Time   `json:"user_updatedat" sql:"user_updatedat"`
	Role       int         `json:"user_role" sql:"user_role"`
	Include    interface{} `json:"include,omitempty"`
}
