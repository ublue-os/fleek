FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"'

RUN ./fleek man > fleek.1 2> /dev/null ; \
  gzip -6 fleek.1 ; \
  exit 0

ENTRYPOINT ["/app/fleek"]
