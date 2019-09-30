START TRANSACTION;

DROP TABLE IF EXISTS versions;

CREATE TABLE  IF NOT EXISTS  institutions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(255) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE submissions ADD institution_id int(11);

ALTER TABLE submissions ADD CONSTRAINT fk_institution
   FOREIGN KEY (institution_id) REFERENCES institutions(id);

insert into institutions (name) values
   ("American Bookbinders Museum"),
   ("Arizona State University"),
   ("Boston Public Library"),
   ("British Library"),
   ("College of the Holy Cross"),
   ("Columbia University"),
   ("Eastern Mennonite University"),
   ("Emory and Henry College"),
   ("Fred W. Smith"),
   ("Georgetown College"),
   ("Georgetown University"),
   ("Grove City College"),
   ("Harvard University"),
   ("Hofstra University"),
   ("Hope College"),
   ("Hunter College"),
   ("Iowa State University"),
   ("James Madison University"),
   ("Lehigh University"),
   ("Long Island University"),
   ("Loyola University Chicago"),
   ("Mary Baldwin University"),
   ("Millsaps College"),
   ("Mount Holyoke College"),
   ("New Entry"),
   ("New York University"),
   ("North Carolina State University"),
   ("Northeastern University"),
   ("Northumbria University"),
   ("Ohio State University"),
   ("Oklahoma State University"),
   ("Open University"),
   ("Rhodes College"),
   ("Roanoke College"),
   ("Rosemont College"),
   ("Ryerson University"),
   ("San Jose State University"),
   ("State University of New York - New Paltz"),
   ("Texas A&M University"),
   ("The University of Virginia"),
   ("Tufts University"),
   ("University of British Columbia"),
   ("University of Calgary"),
   ("University of Connecticut"),
   ("University of Iowa"),
   ("University of Louisville"),
   ("University of Mephis"),
   ("University of Miami"),
   ("University of Michigan"),
   ("University of Missouri"),
   ("University of Montreal"),
   ("University of Nebraska - Lincoln"),
   ("University of Nevada - Las Vegas"),
   ("University of North Carolina - Chapel Hill"),
   ("University of Pennsylvania"),
   ("University of Queensland"),
   ("University of Richmond"),
   ("University of South Carolina"),
   ("University of Tennessee - Chattanooga"),
   ("University of Victoria"),
   ("University of Virginia"),
   ("University of Waterloo"),
   ("Victoria and Albert Museum"),
   ("Virginia Commonwealth University"),
   ("Virginia Military Institute"),
   ("Virginia Polytechnic and State University"),
   ("Virginia Union University"),
   ("Walters Art Museum"),
   ("Washington and Lee University"),
   ("Wellesley College"),
   ("William and Mary University"),
   ("Yale University");

COMMIT;



