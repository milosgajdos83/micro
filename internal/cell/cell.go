package cell

var (
	cells = map[string]string{
		"go": goTmpl,
	}
)

var (
	goTmpl = `
FROM golang:1.13-alpine as build-env
ENV GO111MODULE=on
RUN apk --no-cache add make git gcc libtool musl-dev
WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cell

FROM alpine:latest
RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*
COPY --from=build-env /cell .
ENTRYPOINT ["/cell"]
`
)
