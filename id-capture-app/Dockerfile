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
# Install Tesseract OCR development libraries
RUN apk add --no-cache tesseract-ocr-dev tesseract-ocr
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN go build -o main .

# Final stage
FROM alpine:latest
# Install Tesseract OCR runtime libraries and English language pack
RUN apk --no-cache add ca-certificates tesseract-ocr tesseract-ocr-data-eng
WORKDIR /root/
COPY --from=backend-build /app/main .
COPY --from=frontend-build /app/build ./frontend/build
COPY .env .
EXPOSE 8080
CMD ["./main"]