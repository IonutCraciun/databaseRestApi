servicemachine:
	(cd cmd/servicemachine && \
	go build -o ../../servicemachine)

start: servicemachine
	./servicemachine &

stop:
	kill -9 $(shell ps aux | grep servicemachine | head -1 | cut -d' ' -f2)




