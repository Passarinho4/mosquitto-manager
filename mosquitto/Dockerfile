FROM ubuntu:18.04
RUN \
    apt update && \
    apt install -y mosquitto && \
    rm -rf /var/lib/apt/lists/*

COPY mosquitto.conf /etc/mosquitto/mosquitto.conf
COPY pskfile /etc/mosquitto/pskfile
EXPOSE 8883
ENTRYPOINT ["/usr/sbin/mosquitto","-c", "/etc/mosquitto/mosquitto.conf"]
