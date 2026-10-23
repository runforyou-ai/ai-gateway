ARG TARGETARCH
FROM debian:bookworm-slim
ARG TARGETARCH
RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/* \
	&& mkdir -p /var/log/lite-ai-gateway /etc/lite-ai-gateway
COPY build/lite-ai-gateway-linux-${TARGETARCH} /usr/local/bin/lite-ai-gateway
VOLUME ["/var/log/lite-ai-gateway"]
EXPOSE 3000
CMD ["lite-ai-gateway", "serve", "-f", "/etc/lite-ai-gateway/config.yaml", "--pid", "/tmp/lite-ai-gateway.pid"]
