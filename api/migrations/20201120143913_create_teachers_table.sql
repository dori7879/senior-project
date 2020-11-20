-- +goose U
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `teachers`
(
    `teacher_id`   INT UNSIGNED NOT NULL,
    PRIMARY KEY (`teacher_id`),
    CONSTRAINT `fk_teachers_users` FOREIGN KEY (`teacher_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `teachers`;