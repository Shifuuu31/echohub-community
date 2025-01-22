-- Insert Users
INSERT INTO UserTable (username, email, password)
VALUES 
    ('johndoe', 'john@example.com', 'password123'),
    ('janedoe', 'jane@example.com', 'password456'),
    ('alice', 'alice@example.com', 'password789'),
    ('bobsmith', 'bob@example.com', 'password321'),
    ('charliebrown', 'charlie@example.com', 'password654');

-- Insert Posts
INSERT INTO PostTable (user_id, title, post_content)
VALUES 
    (1, 'AI in 2025: The Future of Machine Learning', 'This post explores the advancements in AI and machine learning in the next decade.'),
    (1, 'AI Ethics and Responsible Use', 'A discussion on the ethical implications of AI and machine learning algorithms.'),
    (1, 'AI and Healthcare: A Revolution in Medicine', 'Exploring the role of AI in revolutionizing the healthcare industry.'),
    (1, 'AI in Education: Personalized Learning', 'How AI is changing the way we approach personalized learning in education.'),
    (1, 'AI in Art: Machines that Create', 'Exploring the role of AI in the world of art creation and design.'),
    
    (2, 'Best Anime of 2025: What to Watch', 'A list of the top anime shows to look forward to in 2025.'),
    (2, 'Exploring the World of Japanese Animation', 'An in-depth look at the history and evolution of Japanese animation.'),
    (2, 'Anime and Its Cultural Impact', 'How anime has influenced global pop culture and entertainment.'),
    (2, 'The Best Anime Movies to Watch in 2025', 'A guide to the most anticipated anime movies of 2025.'),
    (2, 'Anime Streaming Services: Which One to Choose', 'A comparison of the best anime streaming platforms available in 2025.'),
    
    (3, 'Art in the Digital Age: New Media', 'How digital tools and technologies are transforming the world of art.'),
    (3, 'The Intersection of Art and Technology', 'Exploring how artists are using technology to push creative boundaries.'),
    (3, 'Top Art Movements of 2025', 'A look at the emerging art movements that will dominate the art scene in 2025.'),
    (3, 'Art for Social Change', 'How contemporary artists are using their work to promote social and political change.'),
    (3, 'The Evolution of Art in the 21st Century', 'A look at how art has evolved in the 21st century, with a focus on new mediums.'),
    
    (4, 'Exploring the Business of Technology', 'An overview of the tech industry and how it drives the global economy.'),
    (4, 'The Rise of Startups in 2025', 'Why startups are gaining more traction in the business world and what that means for investors.'),
    (4, 'Blockchain Technology: Beyond Cryptocurrencies', 'Exploring the potential applications of blockchain technology outside of cryptocurrency.'),
    (4, 'How Artificial Intelligence is Shaping the Future of Business', 'A look at how AI is being used to optimize business operations and create new opportunities.'),
    (4, 'Business Strategy in 2025: What to Expect', 'Emerging trends in business strategy and how companies are preparing for the future.'),
    
    (5, 'Top Fashion Trends of 2025', 'A rundown of the most popular fashion trends expected in 2025.'),
    (5, 'Sustainable Fashion: The Future of Clothing', 'How the fashion industry is turning towards sustainability and eco-friendly practices.'),
    (5, 'The Impact of Social Media on Fashion Trends', 'How social media platforms like Instagram and TikTok are influencing fashion trends.'),
    (5, 'Fashion Icons to Watch in 2025', 'A look at the fashion influencers and icons who will shape the style landscape in 2025.'),
    (5, 'The Intersection of Technology and Fashion', 'How technology is transforming the fashion industry, from smart clothing to fashion tech startups.');

-- Associate Posts with Categories
INSERT INTO Categories_Posts (category_id, post_id)
VALUES
    -- AI Category (1)
    (1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
    
    -- Anime Category (2)
    (2, 6), (2, 7), (2, 8), (2, 9), (2, 10),
    
    -- Arts Category (3)
    (3, 11), (3, 12), (3, 13), (3, 14), (3, 15),
    
    -- Business Category (4)
    (4, 16), (4, 17), (4, 18), (4, 19), (4, 20);
