CREATE TABLE UserTable(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    creation_date DATE_TIME
);

CREATE TABLE PostTable(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT NOT NULL,
    post_content TEXT NOT NULL,
    categories TEXT NOT NULL,
    creation_date DATE_TIME,
    FOREIGN KEY(user_id) REFERENCES UserTable(id)
);

CREATE TABLE CommentTable(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER,
    comment_content TEXT NOT NULL,
    creation_date DATE_TIME,
    FOREIGN KEY(user_id) REFERENCES UserTable(id),
    FOREIGN KEY(post_id) REFERENCES PostTale(id)
);

CREATE TABLE Likess_Dislike(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    comment_id INTEGER,
    user_id INTEGER,
    liked BOOLEAN,
    FOREIGN KEY(post_id) REFERENCES PostTale(id),
    FOREIGN KEY(comment_id) REFERENCES CommentTable(id),
    FOREIGN KEY(user_id) REFERENCES UserTable(id)
);