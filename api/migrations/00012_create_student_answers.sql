-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `student_answers`
(
    `id`                            INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `type`                          VARCHAR(255) NOT NULL,
    `open_answer`                   TEXT         NULL,
    `true_false_answer`             BOOLEAN      NULL,
    `multiple_choice_answer`        INT UNSIGNED NULL,
    `comments`                      VARCHAR(255) NULL,
    `quiz_submission_id`            INT UNSIGNED NOT NULL,
    `open_question_id`              INT UNSIGNED NULL DEFAULT NULL,
    `true_false_question_id`        INT UNSIGNED NULL DEFAULT NULL,
    `multiple_choice_question_id`   INT UNSIGNED NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`quiz_submission_id`) REFERENCES `quiz_submissions`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
    FOREIGN KEY (`open_question_id`) REFERENCES `open_questions`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
    FOREIGN KEY (`true_false_question_id`) REFERENCES `true_false_questions`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
    FOREIGN KEY (`multiple_choice_question_id`) REFERENCES `multiple_choice_questions`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `student_answers`;