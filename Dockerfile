FROM postgres

WORKDIR /app

COPY . .

CMD ["./main", "main.go"]

