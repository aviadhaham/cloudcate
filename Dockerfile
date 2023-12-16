FROM golang:1.21.4-alpine3.18 AS build
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o main cmd/main.go

FROM golang:1.21.4-alpine3.18
WORKDIR /
COPY --from=build /app/main /
COPY --from=build /app/static /static
ENV PORT=80
EXPOSE $PORT

# Run the executable
CMD ["/main"]