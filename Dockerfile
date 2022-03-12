FROM golang:alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

RUN apk update \
    && apk add --no-cache git ca-certificates tzdata \
    && update-ca-certificates

RUN adduser -D -g '' appuser

ADD . ${GOPATH}/src/app/
WORKDIR ${GOPATH}/src/app

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/hasura_exporter

# --------------------------------------------------------------------------------

FROM gcr.io/distroless/base

ENV WEB_ADDR=9921
ENV DEBUG=false

LABEL summary="Hasura Prometheus exporter" \
      description="A Prometheus exporter for Hasura" \
      name="zolamk/hasura-exporter" \
      url="https://github.com/zolamk/hasura-exporter"

COPY --from=builder /go/bin/hasura_exporter /usr/bin/hasura_exporter
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 9921

ENTRYPOINT ["/usr/bin/hasura_exporter"]
