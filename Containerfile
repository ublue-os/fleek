FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY fleek .
COPY fleek.man.1.gz .

ENTRYPOINT ["/app/fleek"]
