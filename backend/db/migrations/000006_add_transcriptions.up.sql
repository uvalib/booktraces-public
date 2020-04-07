START TRANSACTION;

CREATE TABLE IF NOT EXISTS  transcriptions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   submission_file_id int(11) NOT NULL,
   transcriber_name varchar(255) NOT NULL,
   transcriber_email varchar(255) NOT NULL,
   transcription text NOT NULL,
   submitted_at datetime default CURRENT_TIMESTAMP,
   revision_number SMALLINT NOT NULL default 1,
   approved boolean default false,
   FOREIGN KEY (submission_file_id) REFERENCES submission_files(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

COMMIT;