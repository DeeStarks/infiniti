FROM golang:1.17

RUN apt-get -y update && \
    apt-get -y upgrade

WORKDIR /infiniti
VOLUME ${PWD}:/infiniti

COPY . /infiniti

# The following RUN commands aren't working, it'll be looked into later
# For now, we'll just use the CMD

# RUN go get -d -v ./...
# RUN go install -v ./...

CMD ["/bin/bash", "-c", "go get -d -v ./...; \
    go install -v ./...; \
    go build -o ./bin/infiniti ./cmd/main.go; \
    ./bin/infiniti"]