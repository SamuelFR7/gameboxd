CREATE TABLE IF NOT EXISTS "games" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "slug" TEXT UNIQUE,
    "api_id" INTEGER NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "users" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password_hash" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "sessions" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "user_id" VARCHAR(255) NOT NULL REFERENCES "users" ("id"),
    "expires_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "profiles" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "user_id" VARCHAR(255) REFERENCES "users" ("id"),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "favorite_game_one" VARCHAR(255) REFERENCES "games" ("id"),
    "favorite_game_two" VARCHAR(255) REFERENCES "games" ("id"),
    "favorite_game_three" VARCHAR(255) REFERENCES "games" ("id"),
    "favorite_game_four" VARCHAR(255) REFERENCES "games" ("id")
);

CREATE TABLE IF NOT EXISTS "playing" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "user_id" VARCHAR(255) REFERENCES "users" ("id"),
    "game_id" VARCHAR(255) REFERENCES "games" ("id"),
    "started_at" DATE NOT NULL,
    CONSTRAINT "uq_game_user_playing" UNIQUE ("user_id", "game_id")
);

CREATE TABLE IF NOT EXISTS "logs" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "user_id" VARCHAR(255) REFERENCES "users" ("id"),
    "game_id" VARCHAR(255) REFERENCES "games" ("id"),
    "finished_at" DATE NOT NULL,
    "rate" INTEGER NOT NULL,
    "review" TEXT,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS "ratings" (
    "id" VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "user_id" VARCHAR(255) REFERENCES "users" ("id"),
    "game_id" VARCHAR(255) REFERENCES "games" ("id"),
    "rate" INTEGER NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP,
    CONSTRAINT "uq_game_user_rating" UNIQUE ("user_id", "game_id")
);

CREATE INDEX IF NOT EXISTS "idx_games_slug" ON "games" ("slug");
