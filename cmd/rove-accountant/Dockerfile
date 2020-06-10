FROM golang:latest
LABEL maintainer="Marc Di Luzio <marc.diluzio@gmail.com>"

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o rove-accountant -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=$(git describe --always --long --dirty --tags)'" cmd/rove-accountant/main.go

CMD [ "./rove-accountant" ]

