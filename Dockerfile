FROM golang:1.21.5 as builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /arena-nerf

FROM gcr.io/distroless/static-debian12

COPY --from=builder /arena-nerf /

ENTRYPOINT ["/arena-nerf"]
