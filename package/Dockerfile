#
# build the target application
#
FROM public.ecr.aws/docker/library/golang:1.24.4-alpine3.22 AS builder

# update the packages
RUN apk update && apk upgrade && apk add --no-cache make npm

# see: https://stackoverflow.com/questions/69692842/error-message-error0308010cdigital-envelope-routinesunsupported
ENV NODE_OPTIONS=--openssl-legacy-provider

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY backend ./backend
COPY frontend ./frontend
COPY templates ./templates
RUN make linux-full

#
# build the target container
#
FROM public.ecr.aws/docker/library/alpine:3.22

# update the packages
RUN apk update && apk upgrade && apk add bash curl tzdata ca-certificates msttcorefonts-installer && rm -rf /var/cache/apk/* && update-ms-fonts

# image magick support
RUN apk add fftw-double-libs fontconfig freetype ghostscript ghostscript-fonts lcms2 libbz2 libgcc libgomp libheif libjxl libltdl libraw libx11 libxext libxml2 openjpeg pango tiff zlib libwebpmux libwebpdemux
COPY distro/bin/magick /usr/local/bin
RUN ln -s /usr/local/bin/magick /usr/local/bin/convert && ln -s /usr/local/bin/magick /usr/local/bin/identify
COPY distro/etc/ /usr/local/etc

# Create the run user and group
RUN addgroup --gid 18570 sse && adduser --uid 1984 docker -G sse -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home
ENV APP_HOME=/booktraces-public
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/scripts $APP_HOME/bin $APP_HOME/db
RUN chown -R docker $APP_HOME && chgrp -R sse $APP_HOME

# add the migrate tool
ENV MIGRATE_VERSION=v4.18.1
RUN cd $APP_HOME/bin && curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz && rm LICENSE README.md

# Specify the user
USER docker

# port and run command
EXPOSE 8080
CMD ["scripts/entry.sh"]

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

#
# end of file
#
