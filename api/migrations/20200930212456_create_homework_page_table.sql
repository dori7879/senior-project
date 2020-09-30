-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS homework_page
(
    id             INT UNSIGNED NOT NULL AUTO_INCREMENT,
    title          VARCHAR(255) NOT NULL,
    content        TEXT         NULL,
    created_at     TIMESTAMP    NOT NULL,
    updated_at     TIMESTAMP    NULL,
    locking_at     TIMESTAMP    NULL,
    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS homework_page;