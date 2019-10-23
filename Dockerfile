ARG GO_VERSION=1.13
ARG ALPINE_VERSION=3.10
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

RUN apk --update add \
    gcc \
    libc-dev \
    make


COPY . /go/src/github.com/dhiltgen/site
WORKDIR /go/src/github.com/dhiltgen/site

FROM builder as build
RUN go build -mod=vendor -o /go/bin/site ./cmd

# Unit test setup
FROM builder as test
RUN go test -mod=vendor -cover -covermode=count -coverprofile=/cover.out -v github.com/dhiltgen/site/...
RUN go tool cover -html=./cover.out -o /cover.html
RUN go fmt github.com/dhiltgen/site/...
RUN go list github.com/dhiltgen/site/... | grep -v /vendor/ | xargs golint -set_exit_status

# Target for easily extracing the coverage log
FROM scratch as coverage
COPY --from=test /cover.* /


ARG ALPINE_VERSION=3.10
FROM alpine:${ALPINE_VERSION} as daemon
COPY --from=build /go/src/github.com/dhiltgen/site/templates /templates
COPY --from=build /go/bin/site /bin/site
ENTRYPOINT ["/bin/site"]
