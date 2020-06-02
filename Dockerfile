FROM golang:latest
LABEL maintainer="Marc Di Luzio <marc.diluzio@gmail.com>"

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o rove-server -ldflags="-X version.Version=$(git describe --always --long --dirty)" .

CMD [ "./rove-server" ]

