#
# build the target application
#
FROM golang:1.16.0-alpine3.13 as builder
RUN apk update && apk upgrade && apk add make yarn

# see https://github.com/nodesource/distributions/blob/master/README.md
#RUN apt-get -y install apt-transport-https
#RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
#RUN apt-get -y install nodejs

# see https://yarnpkg.com/en/docs/install#debian-stable
#RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
#RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
#RUN apt-get -y update && apt-get -y install yarn 

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY backend ./backend
COPY frontend ./frontend
COPY templates ./templates
RUN make linux-full

#
# build the target container
#
FROM alpine:3.13

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates imagemagick msttcorefonts-installer && rm -rf /var/cache/apk/* && update-ms-fonts

# Create the run user and group
RUN addgroup --gid 18570 sse && adduser --uid 1984 docker -G sse -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home 
ENV APP_HOME /booktraces-public
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/scripts $APP_HOME/bin $APP_HOME/db
RUN chown -R docker $APP_HOME && chgrp -R sse $APP_HOME

# add the migrate tool
RUN cd $APP_HOME/bin && wget -O - https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz && mv migrate.linux-amd64 migrate

# Specify the user
USER docker

# port and run command
EXPOSE 8080
CMD scripts/entry.sh

# Move in necessary assets
COPY package/scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY backend/db/migrations/*.sql $APP_HOME/db/
COPY package/data/container_bash_profile /home/docker/.profile
COPY --from=builder /build/bin/btsrv.linux $APP_HOME/bin/btsrv
COPY --from=builder /build/bin/public $APP_HOME/bin/public
COPY --from=builder /build/bin/templates $APP_HOME/bin/templates

# Add the build tag
ARG BUILD_TAG
RUN test -n "$BUILD_TAG" && touch $APP_HOME/buildtag.build-$BUILD_TAG || touch $APP_HOME/buildtag.build-0
