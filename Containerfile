FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' github.com/ublue-os/fleek/cmd/fleek


ENTRYPOINT ["/app/fleek"]
