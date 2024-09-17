################
# BUILD BINARY #
################
# golang:1.19-alpine
FROM golang@sha256:46752c2ee3bd8388608e41362964c84f7a6dffe99d86faeddc82d917740c5968 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/go-skeleton
COPY . .

RUN echo $PWD && ls -lah

# Fetch dependencies.
# RUN go get -d -v
RUN go mod download
RUN go mod verify

# CMD go build -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/go-skeleton .

#####################
# MAKE SMALL BINARY #
#####################
FROM alpine:3.14

RUN apk update && apk add --no-cache tzdata
ENV TZ=UTC

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/go-skeleton/resources/templates /go/bin/templates

# Copy the executable.
COPY --from=builder /go/bin/go-skeleton /go/bin/go-skeleton
