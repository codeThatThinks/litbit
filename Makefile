PREFIX='/usr/local'

all: server client

server:
	go build litbit-server.go

client:
	go build litbit.go

install:
	install -D litbit-server ${PREFIX}/bin
	install -D litbit ${PREFIX}/bin
	install -D -m 644 litbit-server.service /usr/lib/systemd/system
	install -D -m 644 litbit.service /usr/lib/systemd/service
	install -D www/* ${PREFIX}/share/litbit/www

clean:
	@rm -f litbit-server
	@rm -f litbit

.PHONY: clean
