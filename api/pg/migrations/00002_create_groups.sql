CREATE TABLE IF NOT EXISTS groups
(
    id                            INT UNSIGNED NOT NULL AUTO_INCREMENT,
    title                         VARCHAR(255) NOT NULL,
    share_link                    VARCHAR(255) NOT NULL,
    owner_id                      INT UNSIGNED NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY (share_link),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
);