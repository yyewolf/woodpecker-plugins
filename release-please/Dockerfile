FROM --platform=$BUILDPLATFORM alpine AS certs

RUN apk --no-cache add ca-certificates
RUN echo "nobody:*:65534:65534:nobody:/_nonexistent:/bin/false" > /passwd
RUN echo "nogroup:*:65534:" >> /group

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT}
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o plugin .

ARG RELEASE_PLEASE_VERSION=v17.1.3
FROM node:20-alpine

LABEL org.opencontainers.image.url="https://github.com/yyewolf/woodpecker-plugins"
LABEL org.opencontainers.image.documentation="https://github.com/yyewolf/woodpecker-plugins/blob/main/release-please/README.md"
LABEL org.opencontainers.image.source="https://github.com/yyewolf/woodpecker-plugins"
LABEL org.goci.wrapper="/release-please"

ARG RELEASE_PLEASE_VERSION

# Install release-please CLI
RUN npm install -g release-please@"$RELEASE_PLEASE_VERSION"

# Copy the plugin binary
COPY --from=builder /app/plugin /release-please
COPY --from=certs /passwd /etc/passwd
COPY --from=certs /group /etc/group

USER 65534:65534

ENTRYPOINT ["/release-please"]
