START TRANSACTION;

alter table events add column event_date2 varchar(50);
update events set event_date2 = date_format( event_date, '%M %e, %Y' );
alter table events drop column event_date;
alter table events change event_date2 event_date varchar(50);

COMMIT;