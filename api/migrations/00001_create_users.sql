CREATE TABLE IF NOT EXISTS users
(
    id             INT UNSIGNED NOT NULL AUTO_INCREMENT,
    first_name     VARCHAR(255) NOT NULL,
    last_name      VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL,
    password_hash  CHAR(60)  NOT NULL,
    date_joined    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login     TIMESTAMP    NULL,
    is_teacher     BOOLEAN      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY (email)
);