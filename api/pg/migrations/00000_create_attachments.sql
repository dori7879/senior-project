CREATE TABLE IF NOT EXISTS attachments
(
    hash          VARCHAR(40) NOT NULL UNIQUE,
    extension       VARCHAR(10) NOT NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    counter       integer NOT NULL DEFAULT 0 CHECK (counter >= 0),
    PRIMARY KEY (hash)
);