FROM cgr.dev/chainguard/go:1.20

WORKDIR /app
# assumes prebuilt binary
COPY fleek .
# assumes prebuilt man page
COPY fleek.1.gz .

ENTRYPOINT ["/app/fleek"]
