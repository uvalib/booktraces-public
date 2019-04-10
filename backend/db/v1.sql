
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
-- Create table for users
--
DROP TABLE IF EXISTS users;
CREATE TABLE users (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   last_name varchar(255) DEFAULT NULL,
   first_name varchar(255) DEFAULT NULL,
   email varchar(255) NOT NULL,
   created_at datetime NOT NULL,
   updated_at datetime NOT NULL,
   UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
