-- May 20, 2016
ALTER TABLE classes ADD CoCode MEDIUMINT NOT NULL DEFAULT 0 AFTER ClassCode;

-- February 14, 2018
-- Add `ImagePath` column to people
ALTER TABLE people ADD COLUMN ImagePath VARCHAR(200);
