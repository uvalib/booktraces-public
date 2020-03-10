START TRANSACTION;

update events set event_date='April 6, 2017' where event_date='April 6-7, 2017';
update events set event_date='November 10, 2016' where event_date='November 10-11, 2016';
update events set event_date='February 4, 2016' where event_date='February 4-5, 2016';
update events set event_date='February 4, 2015' where event_date='February 4 2015';

alter table events add column event_date2 DATE;
update events set event_date2 = str_to_date( event_date, '%M %e, %Y' );
alter table events drop column event_date;
alter table events change event_date2 event_date DATE;

COMMIT;