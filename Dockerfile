############################
# STEP 1 build executable binary
############################
FROM golang as builder


WORKDIR $GOPATH/src/aarnaud/ipxeblue/
COPY . .

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/ipxeblue -mod vendor main.go

############################
# STEP 2 build webui
############################
FROM node:lts-buster as builderui


WORKDIR /webui/
COPY ./webui .

RUN yarn install
RUN yarn build

############################
# STEP 3 ca-certificates
############################
FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates


############################
# STEP 4 build a small image
############################
FROM scratch

ENV GIN_MODE=release
WORKDIR /app/
# Import from builder.
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/ipxeblue /app/ipxeblue
COPY --from=builderui /webui/build /app/ui
COPY templates /app/templates
ENTRYPOINT ["/app/ipxeblue"]
EXPOSE 8080