
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Create table to track schema version
--
DROP TABLE IF EXISTS versions;
CREATE TABLE versions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   version varchar(10) NOT NULL UNIQUE KEY,
   created_at datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
insert into versions(version, created_at) values ("v1", NOW());

--
-- Create table for tags
--
DROP TABLE IF EXISTS tags;
CREATE TABLE tags (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(255) DEFAULT NULL,
   UNIQUE KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add seed tags
insert into tags(name)
   values ("Ownership inscription"), ("Gift Inscription"), ("Date"), ("Place name"), 
   ("Underlining"), ("Annotation"), ("Commentary"), ("Indexing"), ("Drawing"), 
   ("Poem"), ("Quotation"), ("Insert: paper"), ("Insert: botanical"), ("Insert: object");

--
-- Create table for events
--
DROP TABLE IF EXISTS events;
CREATE TABLE events (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   event_date varchar(50)  NOT NULL,
   description varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add seed events
insert into events(event_date,description) values 
("April 16, 2018", "University of Iowa talk and event"),
("March 30, 2018", 'University of Nebraska <a href="https://events.unl.edu/english/2018/03/30/127994">talk and event, </a>with a wonderful news story <a href="https://unlincoln.exposure.co/treasure-hunt">here</a>, another one <a href="https://news.unl.edu/newsrooms/today/article/treasure-hunt-in-library-stacks-uncovers-links-to-university-s-past/">here</a> and an NPR story on the event <a href="http://netnebraska.org/article/news/1123483/project-traces-personal-history-written-old-books">here</a>.'),
("September 11, 2017", '<a href="http://events.louisville.edu/event/andrew_stauffer_garden_of_verses_poems_flowers_and_nineteenth-century_readers">talk</a> and <a href="http://events.louisville.edu/event/traces_in_the_stacks">event </a>at University of Louisville'),
("June 22, 2017", "talk at UNC-Chapel Hill (British Women Writer’s Conference)"),
("April 6-7, 2017", 'talk and event at Holy Cross, with a news story <a href="http://news.holycross.edu/blog/2017/07/13/holy-cross-explores-the-secret-life-of-books/">here</a>.'),
("March 9, 2017", 'talk at Stevenson University, with <a href="http://stevensonenglish.org/eng38101-licastro17/tag/book-traces/">student reactions here</a>'),
("January 26, 2017", "talk at University of Montreal Digital Humanities Colloquium"),
("November 18, 2016", "event at Arizona State University, led by Devoney Looser"),
("November 10-11, 2016", "talk and event at the University of South Carolina"),
("October 26, 2016", "talk and event at Millsaps College"),
("February 4-5, 2016", '<a href="https://www.unlv.edu/event/university-forum-lecture-digital-humanities" target="_blank" rel="noopener"> talk and event the the University of Nevada-Las Vegas</a>; more <a href="https://www.facebook.com/events/162051374169432/" target="_blank" rel="noopener">here</a>.'),
("January 14, 2016", 'talk and event at University of Victoria, with a report on the event <a href="http://uvicspeccoll.tumblr.com/post/138048981427/unravelling-the-codex-dr-andrew-stauffer" target="_blank" rel="noopener">here</a>.'),
("September 24, 2015", ' <a href="http://library.miami.edu/blog/2015/09/04/um-libraries-to-host-hunt-for-hidden-treasures-in-richter-library-stacks/" target="_blank" rel="noopener">Book Traces day at University of Miami</a>, with a report on the event <a href="http://www.themiamihurricane.com/2015/09/27/richter-librarys-book-traces-event-gives-intimate-insight-on-texts-histories/" target="_blank" rel="noopener">here</a>.'),
("May 14, 2015", "talk at Books and/as New Media event, Harvard University"),
("May 5, 2015", "talk at New York University"),
("April 2, 2015", '<a href="http://alumni.harvard.edu/events/consuming-written-word-harvard-library-strategic-conversation"> talk at Harvard University Library</a> (with Robert Darnton and Roger Chartier)'),
("March 25, 2015", 'talk in the Advanced Research Collaborative, CUNY Graduate Center, reviewed by Timothy Griffiths <a href="http://arc.commons.gc.cuny.edu/2015/03/27/the-haptic-lives-and-loves-of-library-books-reflections-on-andrew-stauffers-traces-in-the-stacks-the-troubled-archive-of-nineteenth-century-print/">here</a>'),
("February 4 2015", "talk in the Victorian Seminar, CUNY Graduate Center"),
("February 3, 2015", 'talk in the <a href="http://cencl.uga.edu">18th-and19th Century Colloquium at the University of Georgia</a>'),
("November 21, 2014", 'talk in the <a href="http://depts.washington.edu/text/">Textual Studies Program, University of Washington</a>'),
("October 28, 2014", '<a href="http://library.columbia.edu/locations/rbml/exhibitions/2014-2015.html">talk in the Book History Colloquium at Columbia University</a>, reviewed by Roger Schonfeld <strong><a href="http://www.sr.ithaka.org/blog-individual/notes-columbias-book-history-colloquium">here</a></strong>'),
("October 27, 2014", '<a href="http://www1.lehigh.edu/news/searching-hidden-history-linderman-stacks">talk and event at Lehigh University</a>, reviewed by Annie Johnson <strong><a href="http://anniekjohnson.com/2014/11/03/book-traces-lehigh/">here</a></strong>'),
("October 8, 2014", '<a href="https://blogs.cul.columbia.edu/rbml/2014/08/29/book-traces-comes-to-the-butler-stacks-on-october-8th/"> event at Butler Library, Columbia University</a>, reviewed by Tara Kron<strong> <a href="http://dh.prattsils.org/blog/resources/event-reviews/book-traces-comes-butler-stacks-prof-andrew-stauffer-butler-library-columbia-university-10092014/">here</a></strong><');

--
-- Create table for users
--
DROP TABLE IF EXISTS users;
CREATE TABLE users (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   last_name varchar(255) DEFAULT NULL,
   first_name varchar(255) DEFAULT NULL,
   email varchar(255) NOT NULL,
   token varchar(25) NOT NULL DEFAULT "",
   UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
insert into users (last_name,first_name,email) values ("Foster", "Lou", "lf6f@virginia.edu");

--
-- Create table for submissions
--
DROP TABLE IF EXISTS submissions;
CREATE TABLE submissions (
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

--
-- Create table for submission_files
--
DROP TABLE IF EXISTS submission_files;
CREATE TABLE submission_files (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   submission_id int(11) NOT NULL,
   filename varchar(255) NOT NULL,
   FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for submission_tags
--
DROP TABLE IF EXISTS submission_tags;
CREATE TABLE submission_tags (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   submission_id int(11) NOT NULL,
   tag_id  int(11) NOT NULL,
   FOREIGN KEY (tag_id) REFERENCES tags(id),
   FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
