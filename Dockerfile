FROM golang:latest
LABEL maintainer="Marc Di Luzio <marc.diluzio@gmail.com>"

WORKDIR /app
COPY . .
RUN go mod download

RUN cd cmd/rove-server && go build ./...

CMD "./cmd/rove-server/rove-server"