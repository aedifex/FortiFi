
CREATE DATABASE FortiFi;

USE FortiFi;

Create Table Users (
    id varchar(255) NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(64) NOT NULL,
    email varchar(64) NOT NULL, 
    password varchar(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE RefreshTokens (
    token varchar(255) NOT NULL,
    FK_UserId varchar(255) NOT NULL,
    expires DATETIME NOT NULL,
    PRIMARY KEY (token),
    FOREIGN KEY (FK_UserId) REFERENCES Users(id) ON DELETE CASCADE
);
