ALTER TABLE libraries DROP COLUMN cmd_decider_settings;

ALTER TABLE libraries ADD COLUMN pipeline binary;

ALTER TABLE libraries ADD COLUMN file_cache binary;

DROP TABLE files;

CREATE TABLE IF NOT EXISTS files (
    path text,
    modtime timestamp,
    mediainfo binary
);
