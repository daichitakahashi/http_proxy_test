run:
	go build -o _proxy ./proxy
	HTTPS_PROXY=http://localhost:9999 go run .
