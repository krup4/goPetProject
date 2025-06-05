FROM golang:1.24-alpine3.21 AS builder

RUN apk update && apk add ca-certificates git gcc g++ libc-dev binutils

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o bin/application .

FROM alpine:3.21 AS runner

RUN apk update && apk add ca-certificates libc6-compat openssh bash && rm -rf /var/cache/apk/*

WORKDIR /opt

COPY . . 

COPY --from=builder /opt/bin/application ./

ENV POSTGRES_CONN=postgres://default_user:qwerty123!@db:5432/api_database
ENV POSTGRES_JDBC_URL=jdbc:postgresql://{localhost}:{5432}/{api_database}
ENV POSTGRES_DATABASE=api_database
ENV POSTGRES_USERNAME=default_user
ENV POSTGRES_PASSWORD=qwerty123!
ENV POSTGRES_HOST=db
ENV POSTGRES_PORT=5432

CMD ["bin/application"]