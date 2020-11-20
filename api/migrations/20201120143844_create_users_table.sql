-- +goose U
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `users`
(
    `id`             INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `first_name`     VARCHAR(255) NOT NULL,
    `last_name`      VARCHAR(255) NOT NULL,
    `email`          VARCHAR(255) NOT NULL,
    `password_hash`  CHAR(60)  NOT NULL,
    `date_joined`    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`     TIMESTAMP    NULL,
    `is_active`      BOOLEAN      NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`email`)
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `users`;