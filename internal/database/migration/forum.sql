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
    category_id INTEGER NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE FOREIGN KEY (category_id) REFERENCES Categories(id) ON DELETE CASCADE
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
    category_name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS SessionsUsers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expiration_date DATETIME
);

INSERT
    OR REPLACE INTO Categories (category_name)
VALUES
    ('Movies & Streaming'),
    ('Music & Playlists'),
    ('Gaming & Esports'),
    ('Books & Reads'),
    ('Tech Buzz'),
    ('Gadgets & Gear Reviews'),
    ('Digital Lifestyle'),
    ('Travel Hacks'),
    ('Foodie Finds'),
    ('Health & Wellness'),
    ('DIY Projects & Creatives'),
    ('Skill Up'),
    ('Career Hacks'),
    ('AMA (Ask Me Anything)'),
    ('Meet the Community'),
    ('Life & Relationships Talk'),
    ('Memes & LOLs'),
    ('Global News'),
    ('Local Happenings'),
    ('Hot Takes & Polls'),
    ('Community Challenges'),
    ('Dev & Code Talk'),
    ('AI & The Future'),
    ('Machine Learning Insights'),
    ('Web & App Dev Trends'),
    ('Data Science & Analytics'),
    ('Cybersecurity & Privacy Talk');