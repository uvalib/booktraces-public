START TRANSACTION;

CREATE TABLE IF NOT EXISTS tags (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(255) DEFAULT NULL,
   UNIQUE KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE  IF NOT EXISTS  news (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   title varchar(255)  NOT NULL,
   content text NOT NULL,
   published boolean not null default false,
   created_at datetime NOT NULL default CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE  IF NOT EXISTS  events (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   event_date varchar(50)  NOT NULL,
   description varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS users (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   last_name varchar(255) DEFAULT NULL,
   first_name varchar(255) DEFAULT NULL,
   email varchar(255) NOT NULL,
   token varchar(25) NOT NULL DEFAULT "",
   UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS  submissions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   upload_id varchar(25) DEFAULT NULL,
   submitter_name varchar(255) NOT NULL,
   submitter_email varchar(255) NOT NULL,
   title varchar(255) NOT NULL,
   author varchar(255) NOT NULL,
   publication_info varchar(255) DEFAULT NULL,
   library varchar(255) NOT NULL,
   call_number varchar(20) DEFAULT NULL,
   description text DEFAULT NULL,
   submitted_at datetime default CURRENT_TIMESTAMP,
   public boolean default false,
   UNIQUE KEY (upload_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS  submission_files (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   submission_id int(11) NOT NULL,
   filename varchar(255) NOT NULL,
   FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS  submission_tags (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   submission_id int(11) NOT NULL,
   tag_id  int(11) NOT NULL,
   FOREIGN KEY (tag_id) REFERENCES tags(id),
   FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

COMMIT;
