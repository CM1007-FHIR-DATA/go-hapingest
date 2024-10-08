FROM --platform=$BUILDPLATFORM golang:latest AS builder

WORKDIR /app
COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make


FROM scratch AS runner

COPY --from=builder /app/bin/go-hapingest /app

EXPOSE 8080

ENTRYPOINT [ "/app" ]