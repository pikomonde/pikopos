ServiceName = pikopos

build:
	go build -o service_$(ServiceName)

build-run:
	go build -o service_$(ServiceName) && ./service_$(ServiceName)

stop:
	 - screen -X -S "session_$(ServiceName)" quit
