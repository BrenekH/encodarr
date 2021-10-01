CREATE TABLE IF NOT EXISTS libraries (
    ID integer PRIMARY KEY,
    folder text,
    priority integer,
    fs_check_interval text,
    pipeline binary,
    queue binary,
    file_cache binary,
    path_masks binary
);

CREATE TABLE IF NOT EXISTS files (
    path text UNIQUE,
    modtime timestamp,
    mediainfo binary
);

CREATE TABLE IF NOT EXISTS history (
    time_completed timestamp,
    filename text,
    warnings binary,
    errors binary
);

CREATE TABLE IF NOT EXISTS dispatched_jobs (
    uuid text NOT NULL UNIQUE,
    job binary,
    status binary,
    runner text,
    last_updated timestamp
);
