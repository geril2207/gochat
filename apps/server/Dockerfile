FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ./dist/server ./apps/server

FROM golang:1.22 as runner

WORKDIR /app

COPY --from=builder /app/dist .

RUN ls -la /app

EXPOSE 4001

CMD ["./server"]
