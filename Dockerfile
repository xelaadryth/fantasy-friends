FROM golang:1.9.2-alpine3.7

# Install git
RUN apk update && apk upgrade && apk add --no-cache bash git openssh

# Install app
ADD . /go/src/github.com/xelaadryth/fantasy-friends
WORKDIR /go/src/github.com/xelaadryth/fantasy-friends
RUN go get && go install

CMD go run fantasy-friends.go
