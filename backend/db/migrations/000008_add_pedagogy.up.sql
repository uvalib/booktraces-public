START TRANSACTION;

CREATE TABLE IF NOT EXISTS  pedagogy (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   access_key varchar(100) NOT NULL,
   title varchar(255) NOT NULL,
   content text NOT NULL,
   created_at datetime default CURRENT_TIMESTAMP,
   updated_at datetime default NULL,
   UNIQUE KEY (access_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT into pedagogy (access_key, title, content) VALUES ("index", "Book Traces Pedagogy", "<h1>Placeholder Index Content</h1>");

COMMIT;