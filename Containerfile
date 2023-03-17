FROM cgr.dev/chainguard/go:1.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY build.sh ./

RUN go mod download
RUN go get github.com/ublue-os/fleek/cmd

COPY *.go ./

RUN ./build.sh

CMD [ "/fleek" ]
