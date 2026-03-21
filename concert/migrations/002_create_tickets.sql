CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    concert_id INT NOT NULL,
    status TEXT DEFAULT 'requested'
);