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
# STEP 2 build a small image
############################
FROM scratch

ENV GIN_MODE=release
WORKDIR /app/
# Import from builder.
COPY --from=builder /go/bin/ipxeblue /app/ipxeblue
COPY --from=builderui /webui/build /app/ui
COPY templates /app/templates
ENTRYPOINT ["/app/ipxeblue"]
EXPOSE 8080