# ## Multi-stage build

#
# Init stage, includes logspout source code
# and triggers the build.sh script
#
ARG LOGSPOUT_VERSION=v3.2.13
FROM gliderlabs/logspout:${LOGSPOUT_VERSION} as logspout

#
# Build stage, build logspout with loki adapter
#
FROM golang:1.12.5-alpine3.9 as builder
RUN apk add --update go build-base git mercurial ca-certificates git
ENV GO111MODULE=on
WORKDIR /go/src/github.com/gliderlabs/logspout
COPY --from=logspout /go/src/github.com/gliderlabs/logspout /go/src/github.com/gliderlabs/logspout
COPY modules.go .
COPY ./loki /go/src/github.com/gliderlabs/logspout/adapters/loki
RUN cd /go/src/github.com/gliderlabs/logspout; go mod download
RUN go build -ldflags "-X main.Version=$1" -o /bin/logspout


# #
# # Final stage
# #
FROM alpine
WORKDIR /app
COPY --from=builder /bin/logspout /app/
ENTRYPOINT ["./logspout"]
