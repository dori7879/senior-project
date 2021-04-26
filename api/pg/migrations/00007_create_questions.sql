CREATE TABLE IF NOT EXISTS questions
(
    id                    serial NOT NULL,
    content               TEXT         NOT NULL,
    type                  smallint     NOT NULL,
    fixed                 BOOLEAN      NOT NULL,

    choices               VARCHAR(255)[] NOT NULL,

    open_answer           TEXT         NULL,
    truefalse_answer      BOOLEAN      NULL,
    multiplechoice_answer integer[]        NULL,
    singlechoice_answer   integer          NULL,

    created_at        TIMESTAMP    NOT NULL,
    updated_at        TIMESTAMP    NULL,
    quiz_id           integer  NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE ON UPDATE CASCADE
);