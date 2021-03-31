CREATE TABLE IF NOT EXISTS responses
(
    id                            serial NOT NULL,
    comments                      VARCHAR(255) NOT NULL,
    is_correct                    BOOLEAN      NULL,
    grade                         DECIMAL(3,2) NULL,
    type                          SMALLINT     NOT NULL,

    open_response                 TEXT         NULL,
    truefalse_response            BOOLEAN      NULL,
    multiplechoice_response       integer[]        NULL,
    singlechoice_response         integer          NULL,

    quiz_submission_id            integer  NOT NULL,
    question_id                   integer  NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (quiz_submission_id) REFERENCES quiz_submissions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE ON UPDATE CASCADE
);