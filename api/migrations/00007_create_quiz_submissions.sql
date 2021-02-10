-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `quiz_submissions`
(
    `id`                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `grade`             VARCHAR(255) NOT NULL,
    `comments`          VARCHAR(255) NOT NULL,
    `submitted_at`      TIMESTAMP    NULL,
    `updated_at`        TIMESTAMP    NULL,
    `student_fullname`  VARCHAR(255) NULL,
    `student_id`        INT UNSIGNED NULL DEFAULT NULL,
    `quiz_id`           INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`student_id`) REFERENCES `students`(`student_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`quiz_id`) REFERENCES `quizzes`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `quiz_submissions`;