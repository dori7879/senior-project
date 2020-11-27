-- +goose U
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `homeworks`
(

    `id`                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `content`           TEXT         NULL,
    `grade`             VARCHAR(255) NOT NULL,
    `comments`          VARCHAR(255) NOT NULL,
    `submitted_at`      TIMESTAMP    NOT NULL,
    `updated_at`        TIMESTAMP    NULL,
    `student_fullname`  VARCHAR(255) NULL,
    `student_id`        INT UNSIGNED NULL DEFAULT NULL,
    `homework_page_id`  INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_hw_student` FOREIGN KEY (`student_id`) REFERENCES `students`(`student_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_hw_hwp` FOREIGN KEY (`homework_page_id`) REFERENCES `homework_page`(`id`)
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `homework`;