PREFIX=/usr/local

all: server client

server: server/litbit-server

client: client/litbit

server/litbit-server:
	go build -o server/litbit-server server/*.go

client/litbit:
	go build -o client/litbit client/*.go

install: install-server install-client

install-server:
	install -D server/litbit-server ${PREFIX}/bin/litbit-server
	install -D -m 644 server/litbit-server.service /usr/lib/systemd/system/litbit-server.service
	install -d ${PREFIX}/share/litbit/www
	install server/www/* ${PREFIX}/share/litbit/www

install-client:
	install -D client/litbit ${PREFIX}/bin/litbit
	install -D -m 644 client/litbit.service /usr/lib/systemd/system/litbit.service
	install -d ${PREFIX}/share/litbit
	install client/vlc.sh ${PREFIX}/share/litbit

clean:
	@rm -f server/litbit-server
	@rm -f client/litbit

.PHONY: clean
