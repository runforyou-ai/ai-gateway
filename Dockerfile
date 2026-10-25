ARG TARGETARCH
FROM debian:bookworm-slim
ARG TARGETARCH
RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/* \
	&& mkdir -p /var/log/ai-gateway /etc/ai-gateway
COPY build/ai-gateway-linux-${TARGETARCH} /usr/local/bin/ai-gateway
VOLUME ["/var/log/ai-gateway"]
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
	CMD curl -fsS http://127.0.0.1:3000/health || exit 1
CMD ["ai-gateway", "serve", "-f", "/etc/ai-gateway/config.yaml", "--pid", "/tmp/ai-gateway.pid"]
