FROM golang:1.25.5-alpine3.23 AS builder

ARG GIT_TAG=master

RUN echo GIT_TAG=${GIT_TAG}

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build

######## 
FROM alpine:latest

COPY --from=builder /opt/app/slack-payload-handler /usr/local/bin/slack-payload-handler

CMD ["/usr/local/bin/slack-payload-handler"]
