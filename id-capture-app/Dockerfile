# Build React frontend
FROM node:18-alpine as frontend-build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install --only=production
COPY frontend/ .
RUN npm run build

# Build Go backend
FROM golang:1.21-alpine as backend-build
WORKDIR /app
COPY go.mod ./
RUN go get github.com/gin-gonic/gin github.com/lib/pq github.com/joho/godotenv github.com/google/uuid
RUN go mod tidy
COPY . .
RUN go get .
RUN go build -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-build /app/main .
COPY --from=frontend-build /app/build ./frontend/build
COPY .env .
EXPOSE 8080
CMD ["./main"]