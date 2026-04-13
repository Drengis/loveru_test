FROM golang:1.25-alpine AS builder

ARG APP_DIR
ARG APP_NAME

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ${APP_DIR}/ ./${APP_DIR}/

RUN go build -o /app/${APP_NAME} ./${APP_DIR}

FROM alpine:3.20

ARG APP_NAME
ENV APP_NAME=${APP_NAME}

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/${APP_NAME} /usr/local/bin/app

CMD ["/usr/local/bin/app"]
