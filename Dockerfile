#########################
## G O  B U I L D
#########################
FROM golang:1.18 as build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apt-get update && apt-get install -y make

WORKDIR /app
COPY . .

RUN make bartend

#########################
## W E B
#########################
FROM node:18-alpine as build-web

# home dir is nec. to avoid this issue
#   https://github.com/parcel-bundler/parcel/issues/6578
WORKDIR /home

COPY web/ ./
RUN npm install && npx parcel build --public-url /web/ ./index.html

#########################
## T A R G E T
#########################
FROM alpine:latest

WORKDIR /app

COPY --from=build-web /home/dist/* web/
COPY ./etc/yang /app/yang
COPY ./etc/bartend.json /app/startup.json
COPY --from=build /app/bartend .

EXPOSE 8080
ENV YANGPATH=/app/yang
ENTRYPOINT ["/app/bartend", "-web", "web", "-config", "startup.json"]
