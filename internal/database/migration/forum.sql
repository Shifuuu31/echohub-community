CREATE TABLE IF NOT EXISTS UserTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS PostTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR(70) NOT NULL,
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
    category_name VARCHAR(50) NOT NULL UNIQUE,
    Category_icon_path TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Categories_Posts(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY(category_id) REFERENCES Categories(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES PostTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS SessionsUsers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expiration_date DATETIME,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
);
 
INSERT INTO
    Categories (category_name, Category_icon_path)
VALUES
    (
        'Movies & TV',
        '/assets/imgs/categories-Icons/movies-tv.png'
    ),
    (
        'Music',
        '/assets/imgs/categories-Icons/music.png'
    ),
    (
        'Games',
        '/assets/imgs/categories-Icons/games.png'
    ),
    (
        'Technology',
        '/assets/imgs/categories-Icons/technology.png'
    ),
    ('AI', '/assets/imgs/categories-Icons/ai.png'),
    (
        'Memes',
        '/assets/imgs/categories-Icons/memes.png'
    ),
    (
        'Sports',
        '/assets/imgs/categories-Icons/sports.png'
    ),
    ('News', '/assets/imgs/categories-Icons/news.png'),
    (
        'Fashion',
        '/assets/imgs/categories-Icons/fashion.png'
    ),
    (
        'Science',
        '/assets/imgs/categories-Icons/science.png'
    ),
    ('Art', '/assets/imgs/categories-Icons/art.png'),
    (
        'Anime',
        '/assets/imgs/categories-Icons/anime.png'
    ),
    (
        'Books',
        'assets/imgs/categories-Icons/books.png'
    ),
    (
        'Business',
        '/assets/imgs/categories-Icons/business.png'
    ),
    (
        'Career',
        '/assets/imgs/categories-Icons/career.png'
    ),
    (
        'Culture',
        '/assets/imgs/categories-Icons/culture.png'
    ),
    ('DIY', '/assets/imgs/categories-Icons/diy.png'),
    (
        'Education',
        '/assets/imgs/categories-Icons/education.png'
    ),
    (
        'Podcasts',
        '/assets/imgs/categories-Icons/podcasts.png'
    ),
    (
        'Q&A''s',
        '/assets/imgs/categories-Icons/q-a.png'
    ),
    (
        'Software & Apps',
        '/assets/imgs/categories-Icons/software-apps.png'
    ),
    (
        'Travel',
        '/assets/imgs/categories-Icons/travel.png'
    );
