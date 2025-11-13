FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum / ./
RUN go mod download

COPY ./cmd ./
COPY ./internal ./

RUN CGO_ENABLED=0 GOFLAGS="" go build -trimpath -ldflags "-s -w" -o /out/app ./cmd/main.go


FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/app /app/app

CMD [ "/app/app" ]

