FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY queue.env .
COPY start.sh .
COPY wait-for .
RUN chmod a+x start.sh
RUN chmod a+x wait-for
RUN mkdir logs

EXPOSE 8081
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]