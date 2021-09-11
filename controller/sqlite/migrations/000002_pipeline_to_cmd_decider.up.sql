ALTER TABLE libraries DROP COLUMN pipeline;

ALTER TABLE libraries ADD COLUMN cmd_decider_settings text DEFAULT '';

ALTER TABLE libraries DROP COLUMN file_cache;

DROP TABLE files;

CREATE TABLE IF NOT EXISTS files (
    path text,
    modtime timestamp,
    metadata binary
);
