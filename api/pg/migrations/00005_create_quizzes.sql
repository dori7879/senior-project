CREATE TABLE IF NOT EXISTS quizzes
(
    id                serial NOT NULL,
    title             VARCHAR(255) NOT NULL,
    content           TEXT         NULL,
    max_grade         DECIMAL(5,2) NOT NULL,
    student_link      VARCHAR(255) NOT NULL,
    teacher_link      VARCHAR(255) NOT NULL,
    course_title      VARCHAR(255) NOT NULL,
    mode              VARCHAR(16)  NOT NULL,
    created_at        TIMESTAMP    NOT NULL,
    updated_at        TIMESTAMP    NULL,
    opened_at         TIMESTAMP    NULL,
    closed_at         TIMESTAMP    NULL,
    teacher_fullname  VARCHAR(255) NULL,
    teacher_id        integer  NULL DEFAULT NULL,
    group_id          integer  NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE (student_link, teacher_link),
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE
);