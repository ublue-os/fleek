FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY . ./

RUN ./build.sh

ENTRYPOINT ["/app/fleek"]
