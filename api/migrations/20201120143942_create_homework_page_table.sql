-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `homework_page`
(
    `id`                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`             VARCHAR(255) NOT NULL,
    `content`           TEXT         NULL,
    `student_link`      VARCHAR(255) NOT NULL,
    `teacher_link`      VARCHAR(255) NOT NULL,
    `course_title`      VARCHAR(255) NOT NULL,
    `created_at`        TIMESTAMP    NOT NULL,
    `updated_at`        TIMESTAMP    NULL,
    `opened_at`         TIMESTAMP    NULL,
    `closed_at`         TIMESTAMP    NULL,
    `teacher_fullname`  VARCHAR(255) NULL,
    `teacher_id`        INT UNSIGNED NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_hwp_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `teachers`(`teacher_id`) ON DELETE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `homework_page`;