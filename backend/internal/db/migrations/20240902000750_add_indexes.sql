-- +goose Up
create index idx_users_email on users(email);
create index idx_categories_user_id on goal_categories(user_id);
-- +goose Down
