FROM golang:1.17.2 as build

COPY . /src/weChatRobot-go
WORKDIR /src/weChatRobot-go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags dev -o weChatRobot-go main.go

FROM alpine:3.14

USER root

COPY --from=build /src/weChatRobot-go/weChatRobot-go /weChatRobot-go/

EXPOSE 8080
ENTRYPOINT [ "/weChatRobot-go/weChatRobot-go" ]
