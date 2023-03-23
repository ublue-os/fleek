FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY . ./

RUN ./prep.sh

ENTRYPOINT ["/app/fleek"]
