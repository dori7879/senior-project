-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `open_questions`
(
    `id`                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `content`           TEXT         NOT NULL,
    `answer`            TEXT         NULL,
    `fixed`             BOOLEAN      NOT NULL,
    `created_at`        TIMESTAMP    NOT NULL,
    `updated_at`        TIMESTAMP    NULL,
    `quiz_id`           INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`quiz_id`) REFERENCES `quizzes`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `open_questions`;