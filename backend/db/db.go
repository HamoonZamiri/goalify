package db

import "github.com/jmoiron/sqlx"

var setupScript = `

CREATE TABLE levels if not exists (
    id SERIAL PRIMARY KEY,
    level_up_xp INTEGER,
    cash_reward INTEGER,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE users if not exists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR,
    xp INTEGER DEFAULT 0,
    level_id SERIAL REFERENCES levels(id),
    cash_available INTEGER DEFAULT 0,
    refresh_token UUID DEFAULT gen_random_uuid(),
    refresh_token_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE goal_categories if not exists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR,
    xp_per_goal INTEGER,
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
);`

func New() (*sqlx.DB, error) {
	connStr := "user=postgres dbname=goalify sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.MustExec(setupScript)
	return db, nil
}
