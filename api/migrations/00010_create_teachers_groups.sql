CREATE TABLE IF NOT EXISTS teachers_groups
(
    group_id              INT UNSIGNED NOT NULL,
    teacher_id            INT UNSIGNED NOT NULL,
    PRIMARY KEY (group_id, teacher_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
);