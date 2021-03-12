-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE `quizzes`
ADD `student_group_id` INT UNSIGNED NULL DEFAULT NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE `quizzes`
DROP COLUMN `student_group_id`;