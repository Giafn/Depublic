
ALTER TABLE transactions ADD CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id);
ALTER TABLE transactions ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id);

ALTER TABLE pricings ADD CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id);

ALTER TABLE tickets ADD CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id);
ALTER TABLE tickets ADD CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id);

ALTER TABLE submissions ADD CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id);
ALTER TABLE submissions ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id);

ALTER TABLE profiles ADD CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE;