PREFIX=/usr/local

all: server client

server:
	go build litbit-server.go

client:
	go build litbit.go

install:
	install -D litbit-server ${PREFIX}/bin/litbit-server
	install -D litbit ${PREFIX}/bin/litbit
	install -D -m 644 litbit-server.service /usr/lib/systemd/system/litbit-server.service
	install -D -m 644 litbit.service /usr/lib/systemd/system/litbit.service
	install -d ${PREFIX}/share/litbit/www
	install vlc.sh ${PREFIX}/share/litbit/www
	install www/* ${PREFIX}/share/litbit/www

clean:
	@rm -f litbit-server
	@rm -f litbit

.PHONY: clean
