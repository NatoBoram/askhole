FROM golang:1.24 AS build
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine
WORKDIR /app

COPY --from=build /build/askhole ./
COPY .env .env.local .env.* .env.*.local ./

ENV GO_ENV=production
EXPOSE 9123
CMD ["./askhole"]
