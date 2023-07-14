FROM alpine:latest

RUN apk add -v build-base
RUN apk add -v go 
RUN apk add -v ca-certificates
RUN apk add --no-cache \
    unzip \
    # this is needed only if you want to use scp to copy later your pb_data locally
    openssh

RUN apk add -v nodejs-current npm

# Copy your custom PocketBase and build
COPY . /pb
WORKDIR /pb

# Note: This will pull the latest version of pocketbase. If you are just doing 
# simple customizations and don't have a local build environment for Go, 
# leave this line in. 
# For more complex builds that include other dependencies, remove this 
# line and rely on the go.sum lockfile.

RUN cd client && npm i && PUBLIC_POCKETBASE_URL=https://go-chat-app.fly.dev PUBLIC_WEBSOCKET_URL=wss://go-chat-app.fly.dev/ws npm run build 
RUN go build -tags netgo -ldflags '-s -w' -o output main.go

WORKDIR /

EXPOSE 8080

# start PocketBase
CMD ["./pb/output", "serve", "--http=0.0.0.0:8080"]