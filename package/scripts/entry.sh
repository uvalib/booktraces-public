#!/usr/bin/env bash
#
# run application
#

# run the migrations
echo "Migrate DB"
bin/migrate -path db -verbose -database "mysql://$BOOKTRACES_PUBLIC_DBUSER:$BOOKTRACES_PUBLIC_DBPASS@tcp($BOOKTRACES_PUBLIC_DBHOST:$BOOKTRACES_PUBLIC_DBPORT)/$BOOKTRACES_PUBLIC_DBNAME" up
retVal=$?
if [ $retVal -ne 0 ]; then
   echo "Migration failed"
   exit $?
fi

echo "Launch application"
cd bin; ./btsrv -upload $BOOKTRACES_PUBLIC_UPLOAD_DIR -dbhost $BOOKTRACES_PUBLIC_DBHOST -dbport $BOOKTRACES_PUBLIC_DBPORT -dbname $BOOKTRACES_PUBLIC_DBNAME -dbuser $BOOKTRACES_PUBLIC_DBUSER -dbpass $BOOKTRACES_PUBLIC_DBPASS -smtphost $BOOKTRACES_SMTP_HOST -smtpport $BOOKTRACES_SMTP_PORT -smtpto $BOOKTRACES_EMAIL_TO

#
# end of file
#
