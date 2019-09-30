START TRANSACTION;

ALTER TABLE submissions DROP FOREIGN KEY fk_institution;
ALTER TABLE submissions DROP COLUMN institution_id;
DROP TABLE IF EXISTS institutions;

COMMIT;