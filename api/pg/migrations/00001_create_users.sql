CREATE TABLE IF NOT EXISTS users
(
    id             serial NOT NULL,
    first_name     VARCHAR(255) NOT NULL,
    last_name      VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    password_hash  CHAR(60)  NOT NULL,
    date_joined    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_teacher     BOOLEAN      NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id)
);