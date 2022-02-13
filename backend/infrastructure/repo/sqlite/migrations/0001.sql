CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,

    club TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL,

    auth_token TEXT NOT NULL DEFAULT '',
    auth_time DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00.000'
);
