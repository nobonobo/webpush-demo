FROM golang:1.22.1-alpine3.19 AS build

RUN apk add -u npm git
WORKDIR /root
#RUN git clone https://github.com/nobonobo/webpush-demo
RUN mkdir /root/webpush-demo
COPY ./ /root/webpush-demo/
WORKDIR /root/webpush-demo
RUN npm install && go generate . && go build .

FROM alpine:3.19
COPY --from=build /root/webpush-demo/webpush-demo /app/webpush-demo
WORKDIR /app
ENTRYPOINT [ "/app/webpush-demo" ]
EXPOSE 8080
VOLUME ["/app/subscribes"]
