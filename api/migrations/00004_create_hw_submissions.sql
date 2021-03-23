CREATE TABLE IF NOT EXISTS hw_submissions
(
    id                INT UNSIGNED NOT NULL AUTO_INCREMENT,
    response          TEXT         NULL,
    grade             DECIMAL(3,2) NULL,
    comments          VARCHAR(255) NOT NULL,
    submitted_at      TIMESTAMP    NOT NULL,
    updated_at        TIMESTAMP    NULL,
    student_fullname  VARCHAR(255) NULL,
    student_id        INT UNSIGNED NULL DEFAULT NULL,
    homework_id       INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (homework_id) REFERENCES homeworks(id) ON DELETE CASCADE ON UPDATE CASCADE
);