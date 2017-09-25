# litbit

A DJ for your bluetooth speaker. It plays songs that people request via a web app.

## Installation

Make sure golang is installed and in `$PATH`. Then run the following:

```
make
make install
```

Next start/enable the appropriate systemd service: `litbit-server.service` for the server or `litbit.service` for the client.

## litbit-server.go
Runs on any internet-accesible server and provides the web interface and song queueing.

## litbit.go
Runs on the Raspberry Pi and queries the server for the next song to play. It uses cvlc to play YouTube videos and pianobar to play Pandora stations.

## License

Copyright (c) 2017 Ian Glen

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
