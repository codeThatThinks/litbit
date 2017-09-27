// Client running on raspi

package main

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
	"os/exec"
	"io/ioutil"
	"strings"
	"syscall"
	"os"
	"os/signal"
)

var deviceId string
var shouldQuit = false

func handleExit() {
	channel := make(chan os.Signal, 2)
	signal.Notify(channel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGKILL,
		syscall.SIGSTOP,
	)

	go func() {
		<- channel
		shouldQuit = true
	}()
}

func generateDeviceId() {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, DeviceIdLength)
    for i := 0; i < DeviceIdLength; i++ {
        bytes[i] = byte(rand.Intn(25) + 65)
    }

    deviceId = string(bytes)
}

func main() {
	registered := false

	for !registered {
		// gen unique id
		generateDeviceId()

		// try to register
		resp, err := http.Get(ServerBaseURL + "/" + deviceId + "/register")
		for err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			resp, err = http.Get(ServerBaseURL + "/" + deviceId + "/register")
		}

		if resp.StatusCode == 200 {
			registered = true
		} else if resp.StatusCode != 403 {
			fmt.Printf("Error: Recieved status code %i when trying to register", resp.StatusCode)
			os.Exit(1)
		}
	}

	initLCD()
	drawMessage()
	fmt.Printf("Go to %s/%s to request songs\n", ServerBaseURL, deviceId)
	handleExit()
	
	for !shouldQuit {
		resp, err := http.Get(ServerBaseURL + "/" + deviceId + "/get")
		for err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			resp, err = http.Get(ServerBaseURL + "/" + deviceId + "/get")
		}

		if resp.StatusCode == 200 {
			// play the song
			urlb, _ := ioutil.ReadAll(resp.Body)
			url := strings.TrimSpace(string(urlb))
			fmt.Printf("Playing %s\n", url)
			command := exec.Command("/bin/sh", "vlc.sh", url)
			err = command.Run()
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}
	}

	http.Get(ServerBaseURL + "/" + deviceId + "/unregister")
	fmt.Printf("\n");
	clearLCD()
}
