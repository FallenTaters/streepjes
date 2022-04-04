CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,

    club TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL,

    auth_token TEXT NOT NULL DEFAULT '',
    auth_time DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00.000',

    UNIQUE(username),
    UNIQUE(name)
);

CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,

    UNIQUE(name)
);

CREATE TABLE items (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    price_gladiators INTEGER NOT NULL DEFAULT 0,
    price_parabool INTEGER NOT NULL DEFAULT 0,

    UNIQUE(name),
    FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE members (
    id INTEGER PRIMARY KEY,
    club TEXT NOT NULL,
    name TEXT NOT NULL,
    last_order DATETIME DEFAULT '2000-01-01 00:00:00.000',

    UNIQUE(name)
);

CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    club TEXT NOT NULL,
    bartender_id INTEGER NOT NULL,
    member_id INTEGER,
    contents TEXT,
    price INTEGER NOT NULL,
    order_time DATETIME NOT NULL,
    status TEXT NOT NULL,
    status_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(member_id) REFERENCES members(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY(bartender_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);
