# Build app
FROM golang:1-alpine3.17 AS go-app

WORKDIR /app

COPY . /app

RUN apk add --no-cache gcc alpine-sdk \
    && go mod vendor \
    && CGO_ENABLED=1 CGO_LDFLAGS="-static" go build -o stubman

# =========
FROM alpine:3.17

WORKDIR /app

COPY --from=go-app /app/favicon.* /app/
COPY --from=go-app /app/config.yaml.dist /app/config.yaml
COPY --from=go-app /app/static /app/static
COPY --from=go-app /app/views /app/views
COPY --from=go-app /app/stubman /app/stubman

VOLUME /app

EXPOSE 80

# CMD [ "/app/stubman", "--debug", "-f", "config.yaml"]
ENTRYPOINT [ "/app/stubman" ]
