server:
	go run cmd/main.go

doc:
	echo "Starting swagger generating"
	swag fmt
	swag init -g cmd/main.go --pd
	echo "Swagger generated"