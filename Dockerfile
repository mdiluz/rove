FROM golang:latest
LABEL maintainer="Marc Di Luzio <marc.diluzio@gmail.com>"

WORKDIR /app
COPY . .
RUN go mod download

# Build the executables
RUN go build -o rove -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=$(git describe --always --long --dirty --tags)'" cmd/rove/main.go
RUN go build -o rove-server -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=$(git describe --always --long --dirty --tags)'" cmd/rove-server/main.go

CMD [ "./rove-server" ]

