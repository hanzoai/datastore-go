# hanzoai/datastore-server
# ClickHouse-compatible datastore server
FROM clickhouse/clickhouse-server:25.12-alpine

LABEL org.opencontainers.image.source="https://github.com/hanzoai/datastore-go"
LABEL org.opencontainers.image.description="Hanzo Datastore Server - ClickHouse-compatible columnar database"
LABEL org.opencontainers.image.vendor="Hanzo AI"

# Environment variables
ENV DATASTORE_SKIP_USER_SETUP=${CLICKHOUSE_SKIP_USER_SETUP:-0}

# Copy custom configs if needed
# COPY .docker/datastore/users.xml /etc/clickhouse-server/users.xml

EXPOSE 8123 9000 9009

HEALTHCHECK --interval=10s --timeout=5s --retries=5 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8123/ping || exit 1
