CREATE TABLE wishlists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    event VARCHAR(256) NOT NULL,
    description VARCHAR(1024),
    date TIMESTAMP NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);