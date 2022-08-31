# syntax=docker/dockerfile:experimental

ARG VERSION
ARG REVISION
ARG USER
ARG HOST
ARG BUILD_DATE
ARG BRANCH

FROM golang:1.18.3-alpine AS base
ARG VERSION
ARG REVISION
ARG USER
ARG HOST
ARG BUILD_DATE
ARG BRANCH
ENV VERSION=$VERSION REVISION=$REVISION USER=$USER HOST=$HOST BUILD_DATE=$BUILD_DATE BRANCH=$BRANCH
WORKDIR /src
RUN --mount=type=cache,id=apk,sharing=locked,target=/var/cache/apk ln -vs /var/cache/apk /etc/apk/cache && \
    apk add --update git gcc libc-dev && \
    mkdir /archives && \
    mkdir /dist
COPY . .
WORKDIR /src/lib/web/ui
RUN go generate
WORKDIR /src

FROM base as darwin
RUN GOOS=darwin GARCH=amd64 go build \
            -o /dist/junit-tools \
            -a \
            -ldflags "-s -w -extldflags \"-fno-PIC -static\" -X github.com/kadaan/junit-tools/version.Version=$VERSION -X github.com/kadaan/junit-tools/version.Revision=$REVISION -X github.com/kadaan/junit-tools/version.Branch=$BRANCH -X github.com/kadaan/junit-tools/version.BuildUser=$USER -X github.com/kadaan/junit-tools/version.BuildHost=$HOST -X github.com/kadaan/junit-tools/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo' && \
    tar -czf "/archives/junit_tools_darwin.tar.gz" -C "/dist" .

FROM base as linux
RUN GOOS=linux GARCH=amd64 go build \
            -o /dist/junit-tools \
            -a \
            -ldflags "-s -w -X github.com/kadaan/junit-tools/version.Version=$VERSION -X github.com/kadaan/junit-tools/version.Revision=$REVISION -X github.com/kadaan/junit-tools/version.Branch=$BRANCH -X github.com/kadaan/junit-tools/version.BuildUser=$USER -X github.com/kadaan/junit-tools/version.BuildHost=$HOST -X github.com/kadaan/junit-tools/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo static_build' && \
    tar -czf "/archives/junit_tools_linux.tar.gz" -C "/dist" .

FROM scratch as artifact
COPY --from=darwin /archives/ ./dist/
COPY --from=linux /archives/ ./dist/