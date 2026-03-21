CREATE TABLE IF NOT EXISTS concerts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    tickets_total INT NOT NULL,
    tickets_left INT NOT NULL,
    organizer_email TEXT NOT NULL
);