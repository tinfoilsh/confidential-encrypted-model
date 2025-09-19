FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN go install github.com/modelpack/modctl@v0.1.0-alpha.0
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o manager

FROM quay.io/ifont/skopeo:dev AS skopeo

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/manager .
COPY --from=skopeo /usr/local/bin/skopeo /usr/bin/skopeo
COPY --from=builder /go/bin/modctl /usr/bin/modctl
EXPOSE 8080
CMD ["./manager"]
