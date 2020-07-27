FROM golang:1.14.3 AS builder
ADD /cmd /go/src/mosquitto-manager/cmd
ADD /internal /go/src/mosquitto-manager/internal
ADD go.mod /go/src/mosquitto-manager
WORKDIR /go/src/mosquitto-manager
RUN CGO_ENABLED=0 GOOS=linux go build -a -o mosquitto-manager cmd/main.go

FROM alpine:3.12.0
COPY --from=builder /go/src/mosquitto-manager/mosquitto-manager .
ENTRYPOINT [ "./mosquitto-manager" ]
