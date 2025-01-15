CREATE TABLE IF NOT EXISTS UserTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS PostTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    post_content TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CommentTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    comment_content TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES PostTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Likes_Dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    entity_id INTEGER NOT NULL,
    entity_type TEXT NOT NULL CHECK(entity_type IN ('post', 'comment')),
    liked BOOLEAN NOT NULL,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    categorie_name VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO
    Categories (categorie_name)
VALUES
    ('All') ('My Posts') ('Liked Posts') ('AI') ('Anime') ('Arts') ('Business') ('Career') ('Culture') ('DIY') ('Education') ('Fashion') ('Games') ('Memes') ('Movies & TV') ('Music') ('News') ('Podcasts') ("Q&A's") ('Science') ('Software & Apps') ('Sports') ('Technology') ('Travel');

CREATE TABLE IF NOT EXISTS Categories_Posts(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    categorie_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY(categorie_id) REFERENCES Categories(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES PostTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS SessionsUsers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expiration_date DATETIME,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
);