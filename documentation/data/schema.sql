CREATE DATABASE arcbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE arcbox;

-- Create a `snips` table.
CREATE TABLE snips (
id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
title VARCHAR(100) NOT NULL,
content TEXT NOT NULL,
created DATETIME NOT NULL,
expires DATETIME NOT NULL
);
-- Add an index on the created column.
CREATE INDEX idx_snips_created ON snips(created);