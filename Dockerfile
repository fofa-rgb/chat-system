FROM golang:1.23 AS build_stage
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o app_binary ./cmd/main/main.go
FROM golang:1.22 AS runtime_stage
WORKDIR /app
COPY --from=build_stage /src/app_binary .
EXPOSE 8080
CMD ["./app_binary"]