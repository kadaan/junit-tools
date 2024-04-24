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

FROM base as darwin_intel
RUN GOOS=darwin GARCH=amd64 go build \
            -o /dist/junit_tools_darwin_intel \
            -a \
            -ldflags "-s -w -extldflags \"-fno-PIC -static\" -X github.com/kadaan/junit-tools/version.Version=$VERSION -X github.com/kadaan/junit-tools/version.Revision=$REVISION -X github.com/kadaan/junit-tools/version.Branch=$BRANCH -X github.com/kadaan/junit-tools/version.BuildUser=$USER -X github.com/kadaan/junit-tools/version.BuildHost=$HOST -X github.com/kadaan/junit-tools/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo'
FROM base as darwin_arm
RUN GOOS=darwin GARCH=arm64 go build \
            -o /dist/junit_tools_darwin_arm \
            -a \
            -ldflags "-s -w -extldflags \"-fno-PIC -static\" -X github.com/kadaan/junit-tools/version.Version=$VERSION -X github.com/kadaan/junit-tools/version.Revision=$REVISION -X github.com/kadaan/junit-tools/version.Branch=$BRANCH -X github.com/kadaan/junit-tools/version.BuildUser=$USER -X github.com/kadaan/junit-tools/version.BuildHost=$HOST -X github.com/kadaan/junit-tools/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo'
FROM base as linux
RUN GOOS=linux GARCH=amd64 go build \
            -o /dist/junit_tools_linux \
            -a \
            -ldflags "-d -s -w -extldflags \"-fno-PIC -static\" -X github.com/kadaan/junit-tools/version.Version=$VERSION -X github.com/kadaan/junit-tools/version.Revision=$REVISION -X github.com/kadaan/junit-tools/version.Branch=$BRANCH -X github.com/kadaan/junit-tools/version.BuildUser=$USER -X github.com/kadaan/junit-tools/version.BuildHost=$HOST -X github.com/kadaan/junit-tools/version.BuildDate=$BUILD_DATE" \
            -tags 'osusergo netgo static_build' \
            -installsuffix netgo

FROM scratch as artifact
COPY --from=darwin_intel /dist/ ./dist/
COPY --from=darwin_arm /dist/ ./dist/
COPY --from=linux /dist/ ./dist/