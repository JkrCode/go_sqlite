
-- message schema
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    severity INTEGER NOT NULL,
    descriptionText TEXT NOT NULL,
    receivedDateTime TEXT NOT NULL
);

