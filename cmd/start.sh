source files/etc/env/env.sh

docker-compose up -d

# Do database migration (TODO: Only do synchronous migration with app start in dev)
dbmate up

# Start the HTTP service locally
go run ./cmd/app-http/main.go