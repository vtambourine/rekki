FROM golang:1.14-alpine AS build
WORKDIR /go/src/app/rekki
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./validator ./*.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/rekki/validator /go/bin/
COPY --from=build /go/src/app/rekki/*.txt /go/bin/
EXPOSE 8080
WORKDIR /go/bin
#ENTRYPOINT /go/bin/validator