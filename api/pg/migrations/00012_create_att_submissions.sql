CREATE TABLE IF NOT EXISTS att_submissions
(
    id                serial NOT NULL,
    present           BOOLEAN NOT NULL DEFAULT FALSE,
    submitted_at      TIMESTAMP    NOT NULL,
    updated_at        TIMESTAMP    NULL,
    student_fullname  VARCHAR(255) NULL,
    student_id        integer  NULL DEFAULT NULL,
    attendance_id     integer  NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (attendance_id) REFERENCES attendances(id) ON DELETE CASCADE ON UPDATE CASCADE
);