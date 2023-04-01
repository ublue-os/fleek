FROM scratch
ENTRYPOINT ["/app/fleek"]
COPY fleek /app/fleek
COPY man /
COPY completions /
