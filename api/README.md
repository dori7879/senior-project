##Golang based backend api 

Export variables in .env file into the current shell: `export $(grep -v '^#' .env | xargs -d '\n')`

Run the backend app locally: `go run cmd/app/main.go`

Compile migrate cli app: `go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate`