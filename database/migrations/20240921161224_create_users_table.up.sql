CREATE TABLE users (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password varchar(255) NOT NULL,
    created_at datetime(3) NOT NULL,
    updated_at datetime(3) NOT NULL,
    deleted_at datetime(3) DEFAULT NULL,
    PRIMARY KEY (id),
    KEY idx_users_created_at (created_at),
    KEY idx_users_updated_at (updated_at),
    KEY idx_users_email (email)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
