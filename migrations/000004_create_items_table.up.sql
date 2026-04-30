CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    wishlist_id INTEGER NOT NULL,
    title VARCHAR(256),
    description VARCHAR(1024),
    url VARCHAR(512),
    priority INT NOT NULL,
    is_reserved BOOLEAN NOT NULL DEFAULT FALSE,

    FOREIGN KEY (wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE
);