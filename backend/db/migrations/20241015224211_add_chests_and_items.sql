-- +goose Up
-- +goose StatementBegin
-- Create the enum types first
CREATE TYPE chest_type AS ENUM ('bronze', 'silver', 'gold');
CREATE TYPE item_status AS ENUM ('equipped', 'not_equipped');
CREATE TYPE chest_status AS ENUM ('opened', 'not_opened');
CREATE TYPE item_type AS ENUM ('common', 'rare', 'epic', 'legendary');  -- Assuming this is what `item_type` refers to.

-- Create table `chest_items`
CREATE TABLE chest_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_url VARCHAR,
    title VARCHAR NOT NULL,
    rarity item_type NOT NULL,
    price INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create table `chests`
CREATE TABLE chests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type chest_type NOT NULL,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create table `chest_item_drop_rates`
CREATE TABLE chest_item_drop_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES chest_items(id),
    chest_id UUID REFERENCES chests(id),
    drop_rate FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create table `user_items`
CREATE TABLE user_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    item_id UUID REFERENCES chest_items(id),
    status item_status DEFAULT 'not_equipped',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create table `user_chests`
CREATE TABLE user_chests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    chest_id UUID REFERENCES chests(id),
    quantity_owned INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop tables in reverse order of creation to respect foreign key constraints
DROP TABLE IF EXISTS user_chests;
DROP TABLE IF EXISTS user_items;
DROP TABLE IF EXISTS chest_item_drop_rates;
DROP TABLE IF EXISTS chests;
DROP TABLE IF EXISTS chest_items;

-- Drop enums
DROP TYPE IF EXISTS chest_status;
DROP TYPE IF EXISTS item_status;
DROP TYPE IF EXISTS chest_type;
DROP TYPE IF EXISTS item_type;
-- +goose StatementEnd
