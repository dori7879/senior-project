CREATE TABLE IF NOT EXISTS quiz_submissions
(
    id                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    grade             DECIMAL(3,2) NULL,
    comments          VARCHAR(255) NOT NULL,
    submitted_at      TIMESTAMP    NOT NULL,
    updated_at        TIMESTAMP    NULL,
    student_fullname  VARCHAR(255) NULL,
    student_id        INT UNSIGNED NULL DEFAULT NULL,
    quiz_id           INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE ON UPDATE CASCADE
);