build:
	go build -o service_pikopos

build-run:
	go build -o service_pikopos && ./service_pikopos
