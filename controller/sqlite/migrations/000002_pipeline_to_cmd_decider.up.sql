ALTER TABLE libraries DROP COLUMN pipeline;

ALTER TABLE libraries ADD COLUMN cmd_decider_settings text DEFAULT '';

ALTER TABLE libraries DROP COLUMN file_cache;

DELETE FROM files;

ALTER TABLE files DROP COLUMN mediainfo;

ALTER TABLE files ADD COLUMN metadata binary;
