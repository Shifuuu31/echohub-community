INSERT INTO
    UserTable (username, email, password)
VALUES
    ('user1', 'user1@example.com', 'password1'),
    ('user2', 'user2@example.com', 'password2'),
    ('user3', 'user3@example.com', 'password3'),
    ('user4', 'user4@example.com', 'password4'),
    ('user5', 'user5@example.com', 'password5'),
    ('user6', 'user6@example.com', 'password6'),
    ('user7', 'user7@example.com', 'password7'),
    ('user8', 'user8@example.com', 'password8'),
    ('user9', 'user9@example.com', 'password9'),
    ('user10', 'user10@example.com', 'password10');

INSERT INTO
    PostTable (user_id, title, post_content)
VALUES
    -- Posts for User 1
    (
        1,
        'The Best Movies of 2025',
        'Review of the most popular movies this year.'
    ),
    (
        1,
        'Music Trends in 2025',
        'What’s hot in the music scene right now.'
    ),
    (
        1,
        'Gaming Gadgets Worth Buying',
        'Top gadgets for every gamer this year.'
    ),
    -- Posts for User 2
    (
        2,
        'Top Books to Read in 2025',
        'Incredible reads you should check out.'
    ),
    (
        2,
        'AI Advancements This Year',
        'How AI is shaping the future in 2025.'
    ),
    (
        2,
        'Travel Tips for Summer Adventures',
        'Get the most out of your vacations.'
    ),
    -- Posts for User 3
    (
        3,
        'Best Budget Laptops',
        'Affordable laptops for students and professionals.'
    ),
    (
        3,
        'Health Tips for Gamers',
        'How to stay healthy while gaming.'
    ),
    (
        3,
        'DIY Setup for Home Offices',
        'Creative ways to improve your workspace.'
    ),
    -- Posts for User 4
    (
        4,
        'Movies vs Streaming: What to Choose?',
        'Analyzing the pros and cons.'
    ),
    (
        4,
        'Is Vinyl Making a Comeback?',
        'Exploring the resurgence of vinyl records.'
    ),
    (
        4,
        'Tech Buzzwords Explained',
        'Understanding popular terms in tech.'
    ),
    -- Posts for User 5
    (
        5,
        'The Future of AI',
        'Predictions for AI advancements.'
    ),
    (
        5,
        'Global News Highlights',
        'Major stories from around the world.'
    ),
    (
        5,
        'How to Learn New Skills Fast',
        'Tips for efficient skill development.'
    ),
    -- Posts for User 6
    (
        6,
        'Travel Budget Planning',
        'How to manage your money while traveling.'
    ),
    (
        6,
        'Food Trends of 2025',
        'What’s new in the culinary world.'
    ),
    (
        6,
        'Life Hacks for Busy People',
        'Simplify your daily routine.'
    ),
    -- Posts for User 7
    (
        7,
        'Gadget Reviews: What’s Hot?',
        'Breaking down the latest gadgets.'
    ),
    (
        7,
        'Cybersecurity Basics for Everyone',
        'Protect yourself online.'
    ),
    (
        7,
        'How to Stay Motivated',
        'Tips to keep you going.'
    ),
    -- Posts for User 8
    (
        8,
        'Fitness Routines That Work',
        'Stay fit with these simple exercises.'
    ),
    (
        8,
        'Digital Lifestyle Guide',
        'How to balance online and offline life.'
    ),
    (
        8,
        'Best Local Cafes to Visit',
        'Hidden gems in your city.'
    ),
    -- Posts for User 9
    (
        9,
        'DIY Home Projects',
        'Easy projects for beginners.'
    ),
    (
        9,
        'Skill Up with Online Courses',
        'The best platforms to learn new skills.'
    ),
    (
        9,
        'Career Advice for Developers',
        'Tips for advancing in tech.'
    ),
    -- Posts for User 10
    (
        10,
        'Tech Predictions for 2026',
        'What’s next in the world of technology.'
    ),
    (
        10,
        'Ask Me Anything: Dev Life',
        'Insights into the life of a developer.'
    ),
    (
        10,
        'Hot Takes on Web Development',
        'Opinions and trends in web dev.'
    );

INSERT INTO
    Categories_Posts (categorie_id, post_id)
VALUES
    -- Categories for Posts 1-3 (User 1)
    (1, 1),
    (2, 1),
    (5, 1),
    -- Post 1
    (2, 2),
    (6, 2),
    (7, 2),
    -- Post 2
    (3, 3),
    (5, 3),
    (6, 3),
    -- Post 3
    -- Categories for Posts 4-6 (User 2)
    (4, 4),
    (18, 4),
    (19, 4),
    -- Post 4
    (5, 5),
    (22, 5),
    (23, 5),
    -- Post 5
    (8, 6),
    (9, 6),
    (7, 6),
    -- Post 6
    -- Categories for Posts 7-9 (User 3)
    (6, 7),
    (7, 7),
    (25, 7),
    -- Post 7
    (10, 8),
    (17, 8),
    (16, 8),
    -- Post 8
    (11, 9),
    (12, 9),
    (13, 9),
    -- Post 9
    -- Categories for Posts 10-12 (User 4)
    (1, 10),
    (14, 10),
    (17, 10),
    -- Post 10
    (2, 11),
    (15, 11),
    (18, 11),
    -- Post 11
    (3, 12),
    (16, 12),
    (19, 12),
    -- Post 12
    -- Categories for Posts 13-15 (User 5)
    (22, 13),
    (25, 13),
    (21, 13),
    -- Post 13
    (23, 14),
    (26, 14),
    (20, 14),
    -- Post 14
    (27, 15),
    (25, 15),
    (24, 15),
    -- Post 15
    -- Categories for Posts 16-18 (User 6)
    (8, 16),
    (11, 16),
    (9, 16),
    -- Post 16
    (7, 17),
    (10, 17),
    (6, 17),
    -- Post 17
    (12, 18),
    (19, 18),
    (13, 18),
    -- Post 18
    -- Categories for Posts 19-21 (User 7)
    (5, 19),
    (3, 19),
    (18, 19),
    -- Post 19
    (22, 20),
    (25, 20),
    (27, 20),
    -- Post 20
    (24, 21),
    (26, 21),
    (23, 21),
    -- Post 21
    -- Categories for Posts 22-24 (User 8)
    (10, 22),
    (17, 22),
    (16, 22),
    -- Post 22
    (15, 23),
    (14, 23),
    (13, 23),
    -- Post 23
    (11, 24),
    (12, 24),
    (9, 24),
    -- Post 24
    -- Categories for Posts 25-27 (User 9)
    (2, 25),
    (3, 25),
    (4, 25),
    -- Post 25
    (22, 26),
    (25, 26),
    (24, 26),
    -- Post 26
    (23, 27),
    (21, 27),
    (20, 27),
    -- Post 27
    -- Categories for Posts 28-30 (User 10)
    (7, 28),
    (6, 28),
    (5, 28),
    -- Post 28
    (8, 29),
    (10, 29),
    (9, 29),
    -- Post 29
    (1, 30),
    (2, 30),
    (3, 30);

-- Post 30