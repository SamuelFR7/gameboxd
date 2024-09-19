CREATE TABLE IF NOT EXISTS games (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE,
    api_id INTEGER NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS ratings (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) REFERENCES users(id),
    game_id VARCHAR(255) REFERENCES games(id),
    rate INTEGER NOT NULL
);

ALTER TABLE ratings ADD CONSTRAINT uq_game_user UNIQUE (user_id, game_id); 

CREATE INDEX idx_games_slug ON games (slug);
