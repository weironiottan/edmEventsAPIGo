# Build stage
FROM golang:1.27.2-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/static-debian13:nonroot

# The handlers load templates from ui/html relative to the working directory.
WORKDIR /app
COPY --from=build /out/api ./api
COPY ui ./ui

USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
