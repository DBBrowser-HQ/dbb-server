CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS organizations
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users_organizations
(
    user_id         INTEGER     NOT NULL,
    role            VARCHAR(30) NOT NULL DEFAULT 'observer',
    organization_id INTEGER     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS datasources
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR NOT NULL,
    password        VARCHAR NOT NULL,
    host            VARCHAR NOT NULL,
    port            VARCHAR NOT NULL,
    name            VARCHAR NOT NULL,
    organization_id INTEGER NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
)