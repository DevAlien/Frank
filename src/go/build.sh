go build -o build/mac/frank main.go
env GOOS=linux GOARCH=arm GOARM=5 go build -o build/rpi/frank main.go
