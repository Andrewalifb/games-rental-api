CREATE DATABASE games_rental_api;

-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    address TEXT,
    deposit_amount DECIMAL(10, 2) DEFAULT 0,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Payment History Table
CREATE TABLE payment_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    rental_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'cancelled', 'completed')),
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Game Categories Table
CREATE TABLE game_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Game Platforms Table
CREATE TABLE game_platforms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Games Table
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    availability INT NOT NULL,
    rental_cost DECIMAL(10, 2) NOT NULL,
    platform_id INT REFERENCES game_platforms(id),
    category_id INT REFERENCES game_categories(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Cart Table
CREATE TABLE cart (
    id SERIAL PRIMARY KEY,
    rental_id VARCHAR(255) NOT NULL,
    user_id INT REFERENCES users(id),
    game_id INT REFERENCES games(id),
    quantity INT NOT NULL DEFAULT 1,
    price DECIMAL(10, 2) NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Rental Transactions Table
CREATE TABLE rental_transactions (
    id SERIAL PRIMARY KEY,
    rental_id VARCHAR(255) NOT NULL,
    user_id INT REFERENCES users(id),
    game_id INT REFERENCES games(id),
    payment_id INT REFERENCES payment_history(id),
    quantity INT NOT NULL DEFAULT 1,
    price DECIMAL(10, 2) NOT NULL,
    total_rental_cost DECIMAL(10, 2),
    rented_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    returned_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Rent Maintenance Table
CREATE TABLE rent_maintenance (
    id SERIAL PRIMARY KEY,
    rental_id VARCHAR(255) NOT NULL,
    user_id INT REFERENCES users(id),
    days_left INT,
    status VARCHAR(255) NOT NULL DEFAULT 'not returned' CHECK (status IN ('returned', 'not returned')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert into game_categories
INSERT INTO game_categories (name) VALUES ('Action'), ('Adventure'), ('RPG'), ('Sports'), ('Racing');

-- Insert into game_platforms
INSERT INTO game_platforms (name) VALUES ('PlayStation 5'), ('Xbox One'), ('Nintendo Wii');

-- Insert into games
INSERT INTO games (name, description, availability, rental_cost, platform_id, category_id, created_at, updated_at) VALUES 
('Demons Souls', 'A challenging action RPG game with intense combat and unique multiplayer features.', 10, 100000, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Diablo 4', 'An action RPG game set in a dark fantasy world with cooperative multiplayer.', 10, 110000, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Genshin Impact', 'An open-world action RPG with gacha game mechanics.', 10, 120000, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Final Fantasy 16', 'The latest installment in the long-running RPG series.', 10, 130000, 1, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Control', 'A third-person action-adventure game with a unique physics-based combat system.', 10, 140000, 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Forza Horizon 4', 'A popular open-world racing game with a dynamic weather system and hundreds of cars.', 10, 150000, 2, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Red Dead Redemption 2', 'An open-world action-adventure game set in the late 1800s.', 10, 100000, 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('The Witcher 3: Wild Hunt', 'An open-world action RPG based on the Witcher series of fantasy novels.', 10, 110000, 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Grand Theft Auto V', 'An open-world action-adventure game with a focus on heist missions.', 10, 120000, 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Resident Evil 4', 'A survival horror game with a focus on action and exploration.', 10, 130000, 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('The Legend of Zelda: Skyward Sword', 'An action-adventure game in the Legend of Zelda series with motion controls.', 10, 140000, 3, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Kirby’s Epic Yarn', 'A unique platformer with a charming yarn aesthetic.', 10, 150000, 3, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Kirby’s Return to Dream Land', 'A classic Kirby platformer with cooperative multiplayer.', 10, 100000, 3, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Super Mario Galaxy', 'A 3D platformer set in space with innovative gravity mechanics.', 10, 110000, 3, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Super Smash Bros. Brawl', 'A fighting game featuring characters from various Nintendo franchises.', 10, 120000, 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Super Paper Mario', 'A unique blend of RPG and platforming elements with a charming paper aesthetic.', 10, 130000, 3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Mario Kart Wii', 'A fun and accessible racing game with various Nintendo characters.', 10, 140000, 3, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('WarioWare Smooth Moves', 'A collection of fast-paced microgames with motion controls.', 10, 150000, 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Wii Fit Plus', 'A fitness game with various exercises and balance games.', 10, 100000, 3, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Castlevania The Adventure ReBirth', 'A classic action-adventure game with challenging platforming and combat.', 10, 110000, 3, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

