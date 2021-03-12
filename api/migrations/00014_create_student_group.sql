-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `student_groups`
(
    `id`                            INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`                         VARCHAR(255) NOT NULL,
    `share_link`                    VARCHAR(255) NOT NULL,
    `owner_id`                      INT UNSIGNED NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`owner_id`) REFERENCES `teachers`(`teacher_id`) ON DELETE CASCADE ON UPDATE CASCADE,
) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `student_groups`;