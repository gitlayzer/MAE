FROM golang:1.22.5-alpine as builder
WORKDIR /app
ENV GOPROXY=https://goproxy.io,direct
COPY ./go.mod ./
COPY ./go.sum ./
COPY . .
RUN CGO_ENABLED=0 go build -o mae main.go

FROM busybox as runner
ENV GIN_MODE=release
COPY --from=builder /app/mae /mae
ENTRYPOINT ["/mae", "run"]