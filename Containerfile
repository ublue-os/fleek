FROM cgr.dev/chainguard/go:1.20

WORKDIR /app
# assumes prebuilt binary
COPY fleek .
COPY fleek.1.gz .

ENTRYPOINT ["/app/fleek"]
