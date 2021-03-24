CREATE TABLE IF NOT EXISTS questions
(
    id                    INT UNSIGNED NOT NULL AUTO_INCREMENT,
    content               TEXT         NOT NULL,
    type                  SMALLINT     NOT NULL,
    fixed                 BOOLEAN      NOT NULL,

    choices               VARCHAR(255)[] NOT NULL,

    open_answer           TEXT         NULL,
    truefalse_answer      BOOLEAN      NULL,
    multiplechoice_answer INT[]        NULL,
    singlechoice_answer   INT          NULL,

    created_at        TIMESTAMP    NOT NULL,
    updated_at        TIMESTAMP    NULL,
    quiz_id           INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE ON UPDATE CASCADE
);