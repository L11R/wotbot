-- noinspection SqlResolve

ALTER TABLE stats
    RENAME COLUMN trend_img TO img;
ALTER TABLE stats
    ADD type TEXT DEFAULT 'trend' NOT NULL;
ALTER TABLE stats
    ALTER COLUMN value DROP NOT NULL;