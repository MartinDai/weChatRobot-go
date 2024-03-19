FROM golang:1.21.8 as build

WORKDIR /src/weChatRobot-go

COPY . ./

RUN make

FROM scratch

WORKDIR /app

COPY --from=build /src/weChatRobot-go/bin/weChatRobot_* ./weChatRobot
COPY --from=build /src/weChatRobot-go/config.yml ./

EXPOSE 8080

ENTRYPOINT ["/app/weChatRobot"]
CMD ["-config", "/app/config.yml"]
