FROM golang:1.17.2 as build

COPY . /src/weChatRobot-go
WORKDIR /src/weChatRobot-go
RUN make weChatRobot

FROM alpine:3.15

USER root

COPY --from=build /src/weChatRobot-go/bin/weChatRobot_* /weChatRobot-go/weChatRobot

WORKDIR /weChatRobot-go

EXPOSE 8080
ENTRYPOINT [ "/weChatRobot-go/weChatRobot" ]
