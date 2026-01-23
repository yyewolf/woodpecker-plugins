FROM --platform=\$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=\${TARGETOS} GOARCH=\${TARGETARCH} go build -o plugin .

FROM alpine:latest

RUN apk add --no-cache helm ca-certificates

COPY --from=builder /app/plugin /bin/plugin

ENTRYPOINT ["/bin/plugin"]