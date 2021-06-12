# Multi-Platform Logspout with Loki adapter

Credits for the Loki adapter implementation go to:
https://github.com/darkdarkdragon/logspout-loki

## Build it yourself
Don't trust me? Neither do I!

1. Build it: `docker buildx build --push --pull --platform linux/arm64/v8,linux/amd64 -t <docker-username>/logspout-loki:3.2.11 .`
2. Use it: `docker pull <docker-username>/logspout-loki:3.2.11`

## Prebuilt version
You can also pull from `maronato/logspout-loki`.

Built by Docker-Desktop on MacOS 11.4 using:
```bash
docker buildx build \
  --push --pull \
  --platform linux/arm64/v8,linux/amd64 \
  -t maronato/logspout-loki:3.2.13 \
  -t maronato/logspout-loki:latest .
```
