FROM golang:latest

WORKDIR /app/

ADD . /go/src/github.com/bethanyj28/gomato-slack

RUN go build -o gomato /go/src/github.com/bethanyj28/gomato-slack/cmd/server/main.go
RUN cp /go/src/github.com/bethanyj28/gomato-slack/environment.env .

ENTRYPOINT ["./gomato"]

EXPOSE 8080
