ARG VERSION

FROM --platform=$BUILDPLATFORM golang:1.26.0-alpine@sha256:d4c4845f5d60c6a974c6000ce58ae079328d03ab7f721a0734277e69905473e5 AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -o /build/frontend-service ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot@sha256:a9329520abc449e3b14d5bc3a6ffae065bdde0f02667fa10880c49b35c109fd1 AS runtime
WORKDIR /service
COPY --from=builder /build/frontend-service ./service
COPY internal/frontend/templates /service/frontend/internal/frontend/templates
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8081
ENTRYPOINT ["./service", "-config", "/config/config.yaml"]
