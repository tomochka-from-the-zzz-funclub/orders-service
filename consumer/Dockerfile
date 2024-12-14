
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd
RUN go mod tidy
RUN go build -o api-gateway
FROM golang:1.23
COPY --from=builder /app/cmd/api-gateway /api-gateway
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh


COPY template/getUser.html /app 
COPY template/loginUser.html /app 
COPY template/registrationUser.html /app 
COPY template/loginAdmin.html /app 
COPY template/workAdmin.html /app 

CMD ["/app/wait-for-it.sh", "kafka:9092", "--", "/api-gateway"]
