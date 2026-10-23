ARG TARGETARCH
FROM debian:bookworm-slim
ARG TARGETARCH
RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/* \
	&& mkdir -p /var/log/lite-ai-gateway /etc/lite-ai-gateway
COPY build/lite-ai-gateway-linux-${TARGETARCH} /usr/local/bin/lite-ai-gateway
VOLUME ["/var/log/lite-ai-gateway"]
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
	CMD curl -fsS http://127.0.0.1:3000/health || exit 1
CMD ["lite-ai-gateway", "serve", "-f", "/etc/lite-ai-gateway/config.yaml", "--pid", "/tmp/lite-ai-gateway.pid"]
