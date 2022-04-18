build:
	env GOOS=linux GOARCH=arm GOARM=6 go build -o bin/pumpit_ARMv6 src/*
	go build -o bin/pumpit_x86 src/*

run:
	go run src/*