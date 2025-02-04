-- enable foreign_keys
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS UserTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS UserSessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER UNIQUE NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expiration_date DATETIME,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
);
    
CREATE TABLE IF NOT EXISTS PostTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR(70) NOT NULL,
    content TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
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

CREATE TABLE IF NOT EXISTS CommentTable (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        comment_content TEXT NOT NULL,
        creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES UserTable (id) ON DELETE CASCADE,
        FOREIGN KEY (post_id) REFERENCES PostTable (id) ON DELETE CASCADE
);


-- Uncomment those fields just in need
-- INSERT INTO
--     Categories (category_name, Category_icon_path)
-- VALUES
--     (
--         'Movies & TV',
--         '/assets/imgs/categories-Icons/movies-tv.png'
--     ),
--     (
--         'Music',
--         '/assets/imgs/categories-Icons/music.png'
--     ),
--     (
--         'Games',
--         '/assets/imgs/categories-Icons/games.png'
--     ),
--     (
--         'Technology',
--         '/assets/imgs/categories-Icons/technology.png'
--     ),
--     ('AI', '/assets/imgs/categories-Icons/ai.png'),
--     (
--         'Memes',
--         '/assets/imgs/categories-Icons/memes.png'
--     ),
--     (
--         'Sports',
--         '/assets/imgs/categories-Icons/sports.png'
--     ),
--     ('News', '/assets/imgs/categories-Icons/news.png'),
--     (
--         'Fashion',
--         '/assets/imgs/categories-Icons/fashion.png'
--     ),
--     (
--         'Science',
--         '/assets/imgs/categories-Icons/science.png'
--     ),
--     ('Art', '/assets/imgs/categories-Icons/art.png'),
--     (
--         'Anime',
--         '/assets/imgs/categories-Icons/anime.png'
--     ),
--     (
--         'Books',
--         'assets/imgs/categories-Icons/books.png'
--     ),
--     (
--         'Business',
--         '/assets/imgs/categories-Icons/business.png'
--     ),
--     (
--         'Career',
--         '/assets/imgs/categories-Icons/career.png'
--     ),
--     (
--         'Culture',
--         '/assets/imgs/categories-Icons/culture.png'
--     ),
--     ('DIY', '/assets/imgs/categories-Icons/diy.png'),
--     (
--         'Education',
--         '/assets/imgs/categories-Icons/education.png'
--     ),
--     (
--         'Podcasts',
--         '/assets/imgs/categories-Icons/podcasts.png'
--     ),
--     (
--         'Q&A''s',
--         '/assets/imgs/categories-Icons/q-a.png'
--     ),
--     (
--         'Software & Apps',
--         '/assets/imgs/categories-Icons/software.png'
--     ),
--     (
--         'Travel',
--         '/assets/imgs/categories-Icons/travel.png'
--     );



INSERT INTO PostTable (user_id, title, content, creation_date) VALUES
(1, 'Sample Title 1', 'This is a longer sample content for post 1. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 2', 'This is a longer sample content for post 2. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 3', 'This is a longer sample content for post 3. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 4', 'This is a longer sample content for post 4. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 5', 'This is a longer sample content for post 5. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 6', 'This is a longer sample content for post 6. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 7', 'This is a longer sample content for post 7. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(1, 'Sample Title 8', 'This is a longer sample content for post 8. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 9', 'This is a longer sample content for post 9. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 10', 'This is a longer sample content for post 10. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 11', 'This is a longer sample content for post 11. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 12', 'This is a longer sample content for post 12. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 13', 'This is a longer sample content for post 13. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 14', 'This is a longer sample content for post 14. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 15', 'This is a longer sample content for post 15. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 16', 'This is a longer sample content for post 16. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 17', 'This is a longer sample content for post 17. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 18', 'This is a longer sample content for post 18. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 19', 'This is a longer sample content for post 19. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(2, 'Sample Title 20', 'This is a longer sample content for post 20. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 21', 'This is a longer sample content for post 21. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 22', 'This is a longer sample content for post 22. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 23', 'This is a longer sample content for post 23. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 24', 'This is a longer sample content for post 24. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 25', 'This is a longer sample content for post 25. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 26', 'This is a longer sample content for post 26. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 27', 'This is a longer sample content for post 27. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(3, 'Sample Title 28', 'This is a longer sample content for post 28. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 29', 'This is a longer sample content for post 29. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 30', 'This is a longer sample content for post 30. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 31', 'This is a longer sample content for post 31. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 32', 'This is a longer sample content for post 32. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 33', 'This is a longer sample content for post 33. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 34', 'This is a longer sample content for post 34. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(4, 'Sample Title 35', 'This is a longer sample content for post 35. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(5, 'Sample Title 36', 'This is a longer sample content for post 36. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(5, 'Sample Title 37', 'This is a longer sample content for post 37. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(5, 'Sample Title 38', 'This is a longer sample content for post 38. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(5, 'Sample Title 39', 'This is a longer sample content for post 39. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP),
(5, 'Sample Title 40', 'This is a longer sample content for post 40. It contains more details and explanations about the topic.', CURRENT_TIMESTAMP);
