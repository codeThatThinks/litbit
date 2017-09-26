#!/bin/sh

vlc -I dummy --no-video --network-caching=3000 --disc-caching=3000 --preferred-resolution=360 --play-and-exit "$1"
