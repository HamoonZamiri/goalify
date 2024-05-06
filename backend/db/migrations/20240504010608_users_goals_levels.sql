-- +goose Up
CREATE TABLE levels  (
    id SERIAL PRIMARY KEY,
    level_up_xp INTEGER NOT NULL,
    cash_reward INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE users  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    xp INTEGER DEFAULT 0,
    level_id SERIAL REFERENCES levels(id),
    cash_available INTEGER DEFAULT 0,
    refresh_token UUID DEFAULT gen_random_uuid(),
    refresh_token_expiry TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE goal_categories  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    xp_per_goal INTEGER NOT NULL,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TYPE goal_status AS ENUM ('complete', 'not_complete');

CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    description VARCHAR DEFAULT '',
    user_id UUID REFERENCES users(id),
    category_id UUID REFERENCES goal_categories(id),
    status goal_status DEFAULT 'not_complete',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Insert Default Levels
INSERT INTO levels (id, level_up_xp, cash_reward) VALUES (1, 100, 10);

-- +goose Down
DROP TABLE goals;
DROP TABLE goal_categories;
DROP TABLE users;
DROP TABLE levels;
DROP TYPE goal_status;
