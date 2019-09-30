# Book Traces Public Submission Service

An online form form the public submission of Book Traces Materials 

### System Requirements
* GO version 1.12 or greater (mod required)
* Node version 8.11.1 or higher (https://nodejs.org/en/)
* Yarn version 1.13 or greater
* vue-cli 3 version 3.5 or greater

### Build Instructions

1. After clone, `cd frontend` and execute `yarn install` to prepare the front end
2. Run the Makefile target `build all` to generate binaries for linux, darwin and the front end.  All results will be in the bin directory.

### Database Notes

The backend of the client uses a Postgres DB for user settings/preferences. The schema is managed by 
https://github.com/golang-migrate/migrate and the scripts are in ./backend/db/migrations.

Install the migrate binary on your host system. For OSX, the easiest method is brew. Execute:

`brew install golang-migrate`.

Define your MySQL connection params in an environment variabddle, like this:

`export BTDB=mysql://user:password@tcp(host:port)/booktraces_public`

Run migrations like this:

`migrate -database "${BTDB}" -path backend/db/migrations up`

Example migrate commads to create a migration and run one:

* `migrate create -ext sql -dir backend/db/migrations -seq add_new_table`
* `migrate -database "${BTDB}" -path backend/db/migrations/ up`
