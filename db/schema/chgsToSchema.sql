-- May 20, 2016
ALTER TABLE classes ADD CoCode MEDIUMINT NOT NULL DEFAULT 0 AFTER ClassCode;

-- February 14, 2018
-- Add `ImagePath` column to people
ALTER TABLE people ADD COLUMN ImagePath VARCHAR(200);

-- February 16, 2018
-- Update `ImagePath` column's value for existing user who have profile image
update people set ImagePath="27.jpg" WHERE UID=27;
update people set ImagePath="85" WHERE UID=85;
update people set ImagePath="86.jpg" WHERE UID=86;
update people set ImagePath="103.jpg" WHERE UID=103;
update people set ImagePath="198.jpg" WHERE UID=198;
update people set ImagePath="200.jpg" WHERE UID=200;
update people set ImagePath="202.png" WHERE UID=202;
update people set ImagePath="203.jpg" WHERE UID=203;
update people set ImagePath="207.jpg" WHERE UID=207;
update people set ImagePath="210.jpg" WHERE UID=210;
update people set ImagePath="211.png" WHERE UID=211;
update people set ImagePath="263.jpg" WHERE UID=263;
update people set ImagePath="264.jpg" WHERE UID=264;
update people set ImagePath="267.jpg" WHERE UID=267;
update people set ImagePath="281.jpg" WHERE UID=281;

-- Mar 04, 2018
-- Add UserAgent, IP to sessions table
ALTER TABLE sessions ADD COLUMN UserAgent VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE sessions ADD COLUMN IP VARCHAR(40) NOT NULL DEFAULT '';
