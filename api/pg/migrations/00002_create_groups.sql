CREATE TABLE IF NOT EXISTS groups
(
    id                            serial NOT NULL,
    title                         VARCHAR(255) NOT NULL,
    share_link                    VARCHAR(255) NOT NULL UNIQUE,
    owner_id                      integer NULL DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);