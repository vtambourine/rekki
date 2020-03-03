FROM golang:1.14-alpine AS build
WORKDIR /go/src/app/rekki
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./app .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /go/bin
COPY --from=build /go/src/app/rekki/app /go/bin/
COPY --from=build /go/src/app/rekki/*.txt /go/bin/
EXPOSE 8080
ENTRYPOINT /go/bin/app