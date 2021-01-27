-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `answer_choices`
(
    `id`                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `content`           TEXT         NOT NULL,
    `correct_answer`    BOOLEAN      NOT NULL,
    `question_id`       INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`question_id`) REFERENCES `multiple_choice_questions`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `answer_choices`;