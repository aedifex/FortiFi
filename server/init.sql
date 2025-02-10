
CREATE DATABASE FortiFi;

USE FortiFi;

Create Table Users (
    id varchar(255) NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(64) NOT NULL,
    email varchar(64) NOT NULL,
    fcm_token varchar(255),
    password varchar(255) NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX Users_Email_Index ON Users(email ASC);


CREATE TABLE PiRefreshTokens (
    token_hash varchar(255) NOT NULL,
    id varchar(255) NOT NULL,
    expires DATETIME NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX Pi_Expires_Index ON PiRefreshTokens(expires ASC);

CREATE TABLE UserRefreshTokens (
    token_hash varchar(255) NOT NULL,
    id varchar(255) NOT NULL,
    expires DATETIME NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES Users(id) ON DELETE CASCADE
);
CREATE INDEX Users_Expires_Index ON UserRefreshTokens(expires ASC);

CREATE TABLE NetworkEvents (
    id varchar(255) NOT NULL,
    details varchar(255) NOT NULL,
    ts DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    FOREIGN KEY (id) REFERENCES Users(id)
)