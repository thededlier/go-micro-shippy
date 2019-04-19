FROM alpine:latest

RUN mkdir -p /app
WORKDIR /app

ADD sample-consignment.json /app/sample-consignment.json
ADD consignment-cli /app/consignment-cli

CMD ["./consignment-cli"]
