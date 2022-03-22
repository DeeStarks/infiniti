FROM golang:1.17

RUN apt-get -y update && \
    apt-get -y upgrade
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go get -u github.com/cosmtrek/air
    
COPY . /infiniti

WORKDIR /infiniti
VOLUME ${PWD}:/infiniti

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o ./bin/infiniti ./cmd/main.go

CMD ["air"]