CREATE TABLE IF NOT EXISTS students_groups
(
    group_id              INT UNSIGNED NOT NULL,
    student_id            INT UNSIGNED NOT NULL,
    PRIMARY KEY (group_id, student_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
);