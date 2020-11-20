-- +goose U
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `students`
(
    `student_id`   INT UNSIGNED NOT NULL,
    PRIMARY KEY (`student_id`),
    CONSTRAINT `fk_students_users` FOREIGN KEY (`student_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `students`;