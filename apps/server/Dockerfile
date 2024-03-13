FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ./dist/server ./apps/server

FROM golang:1.22 as runner

ARG SERVER_INTERNAL_PORT

WORKDIR /app

COPY --from=builder /app/dist .

RUN ls -la /app

EXPOSE ${SERVER_INTERNAL_PORT}

CMD ["./server"]
