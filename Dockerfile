#################################################
# Frontend builder
#################################################
FROM node:21-slim AS frontend-builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app

COPY web/package*.json .
COPY web/pnpm-lock.yaml .
RUN pnpm install --frozen-lockfile

COPY web .
RUN pnpm run build

#################################################
# Server builder
#################################################
FROM golang:1.21.4-alpine3.18 AS server-builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o server cmd/main.go

#################################################
# Final image
#################################################

FROM alpine:3.18

WORKDIR /app

COPY --from=frontend-builder /app/dist ./web/dist
COPY --from=server-builder /app/server .

ENV PORT=80
EXPOSE $PORT

# Run the executable
CMD ["./server"]
