# Build stage
FROM golang:1.22 AS builder

WORKDIR /src

# Install protoc + deps
RUN apt-get update && apt-get install -y --no-install-recommends unzip && rm -rf /var/lib/apt/lists/*

# Get googleapis for annotations
RUN git clone --depth 1 https://github.com/googleapis/googleapis /googleapis

COPY go.mod .
RUN go mod download

COPY . .

# Generate descriptor for Envoy transcoder
RUN --mount=type=cache,target=/root/.cache \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

RUN curl -L -o /usr/local/bin/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v27.1/protoc-27.1-linux-x86_64.zip && \
    unzip -o /usr/local/bin/protoc.zip -d /usr/local && rm /usr/local/bin/protoc.zip

RUN mkdir -p /out gen/orders/v1 && \
    protoc -Iproto -I/googleapis \
      --go_out=gen --go-grpc_out=gen \
      --descriptor_set_out=/out/api_descriptor.pb --include_imports \
      proto/order.proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /
COPY --from=builder /out/app /app
COPY --from=builder /out/api_descriptor.pb /api_descriptor.pb

#ENV GRPC_ADDR=:50051 \
#    PG_ADDR=postgres:5432 \
#    PG_USER=postgres \
#    PG_PASS=postgres \
#    PG_DB=ordersdb \
#    JWT_SECRET=technonext_secret

EXPOSE 50051
ENTRYPOINT ["/app"]
