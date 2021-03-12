-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `teachers_student_groups`
(
    `student_group_id`              INT UNSIGNED NOT NULL,
    `teacher_id`                    INT UNSIGNED NOT NULL,
    PRIMARY KEY (`student_group_id`, `teacher_id`),
    FOREIGN KEY (`student_group_id`) REFERENCES `student_groups`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`teacher_id`) REFERENCES `teachers`(`teacher_id`) ON DELETE CASCADE ON UPDATE CASCADE,
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `teachers_student_groups`;