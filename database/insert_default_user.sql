--
-- Insert default admin user
--

INSERT INTO  `users`(`user_id`, `user_name`, `user_email`, `user_password`, `user_tokenhash`, `user_role`)
        VALUES (0, 'Administrator', 'admin@email.com', 'password', 'passwordhash', 42);
