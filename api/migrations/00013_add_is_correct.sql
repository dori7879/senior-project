-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE `student_answers`
ADD `is_correct` BOOLEAN NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE `student_answers`
DROP COLUMN `is_correct`;