# ----- #
# BUILD #
# ----- #

FROM golang:1.27.1-alpine3.24 AS build

# set up workdir
WORKDIR /go/src

# install dependences
COPY ./go.mod .
COPY ./go.sum .
COPY ./vendor ./vendor

# compile app
COPY ./internal ./internal
COPY ./config ./config
COPY ./main.go ./main.go
RUN go build -mod=vendor ./main.go

# --- #
# RUN #
# --- #

FROM alpine:3.21 AS run

WORKDIR /app

# copy compiled app
COPY --from=build /go/src/main .

# run app
CMD ["/bin/sh", "-c", "/app/main"]
