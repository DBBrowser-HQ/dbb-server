CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL
);

CREATE TABLE organizations
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE users_organizations
(
    user_id         INTEGER     NOT NULL,
    role            VARCHAR(30) NOT NULL DEFAULT 'observer',
    organization_id INTEGER     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
);

CREATE TABLE data_sources
(
    id              SERIAL PRIMARY KEY,
    type            VARCHAR NOT NULL,
    username        VARCHAR NOT NULL,
    password        VARCHAR NOT NULL,
    host            VARCHAR NOT NULL,
    port            VARCHAR NOT NULL,
    name            VARCHAR NOT NULL,
    organization_id INTEGER NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
)