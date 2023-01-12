# Build part

FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY app /app
RUN go mod download
COPY healthcheck /healthcheck

RUN CGO_ENABLED=0 go build -v -ldflags "-s -w" -o /health /healthcheck/main.go
RUN CGO_ENABLED=0 go build -o /url-collector

# Exec part
FROM gcr.io/distroless/base-debian10

WORKDIR /

#USER nonroot:nonroot

#COPY --from=builder --chown=nonroot:nonroot --chmod=500 /health /health
#COPY --from=builder --chown=nonroot:nonroot --chmod=500 /url-collector /url-collector
COPY --from=builder --chmod=500 /health /usr/bin/health
COPY --from=builder --chmod=500 /url-collector /usr/bin/url-collector

ARG URL_COLLECTOR_PORT
EXPOSE ${URL_COLLECTOR_PORT}

# HEALTHCHECK
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 CMD [ "health" ]

ENTRYPOINT [ "url-collector" ]
