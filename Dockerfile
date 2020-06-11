FROM golang:latest
LABEL maintainer="Marc Di Luzio <marc.diluzio@gmail.com>"

WORKDIR /app
COPY . .
RUN go mod download

# For /usr/share/dict/words
RUN apt-get -q update && apt-get -qy install wamerican

# Build both executables
RUN go build -o rove-server -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=$(git describe --always --long --dirty --tags)'" cmd/rove-server/main.go
RUN go build -o rove-accountant -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=$(git describe --always --long --dirty --tags)'" cmd/rove-accountant/main.go

CMD [ "./rove-server" ]

