CREATE TABLE IF NOT EXISTS UserTable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS SessionsUsers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expiration_date DATETIME,
    FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
);

INSERT INTO UserTable (username, email, hashed_password) 
VALUES
    ('john_doe', 'john.doe@example.com', 'hashed_password123'),
    ('jane_smith', 'jane.smith@example.com', 'securePass2024'),
    ('alice_wonder', 'alice.wonder@example.com', 'aliceP@ssword'),
    ('bob_marley', 'bob.marley@example.com', 'marleyhashed_password!23'),
    ('charlie_brown', 'charlie.brown@example.com', 'charlieBrown123'),
    ('emily_davis', 'emily.davis@example.com', 'emilyPass#2024'),
    ('daniel_lee', 'daniel.lee@example.com', 'danielPass2024!'),
    ('lucy_harper', 'lucy.harper@example.com', 'lucyH@rper123'),
    ('michael_clark', 'michael.clark@example.com', 'michaelClark!2024'),
    ('olivia_wilson', 'olivia.wilson@example.com', 'oliviaPassword123');
    -- ('shifuuu', 'shifuuu@example.com', '$2a$12$yHoYoooqnJA6ney0kv1GpuTJcQOKyZjESJ01u.VHIPn0gv/idtWY6');
    -- ('abdo', 'abdo@gmail.com', '$2a$12$o7vY6IbHhnVTjkKOmdnY7ON0car9ImXY2ZZHLqc1dePaXHmR.mTsi');
    


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


