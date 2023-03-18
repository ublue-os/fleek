FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY build.sh ./

RUN go mod download

COPY *.go ./
COPY cmd ./cmd
COPY core ./core
COPY locales ./locales
RUN ./build.sh
RUN ./fleek man > fleek.man.1
RUN gzip fleek.man.1
ENTRYPOINT ["/app/fleek"]
