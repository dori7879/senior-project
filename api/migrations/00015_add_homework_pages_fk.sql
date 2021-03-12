-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE `homework_pages`
ADD `student_group_id` INT UNSIGNED NULL DEFAULT NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE `homework_pages`
DROP COLUMN `student_group_id`;