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
)

var id_length = 4
var server_base_url = "http://litbit.in"

func GenerateDeviceId() string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, id_length)
    for i := 0; i < id_length; i++ {
        bytes[i] = byte(rand.Intn(25) + 65)
    }

    return string(bytes)
}

func main() {
	registered := false

	device_id := ""

	for !registered {
		// gen unique id
		device_id = GenerateDeviceId()

		// try to register
		resp, err := http.Get(server_base_url + "/" + device_id + "/register")
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == 200 {
			registered = true
		}
	}

	fmt.Printf("Go to %s/%s to request songs\n", server_base_url, device_id)

	for true {
		resp, err := http.Get(server_base_url + "/" + device_id + "/get")
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == 200 {
			// play the song
			urlb, _ := ioutil.ReadAll(resp.Body)
			url := strings.TrimSpace(string(urlb))
			fmt.Printf("Playing %s\n", url)
			command := exec.Command("/bin/sh", "vlc.sh", url)
			err := command.Run()
			if err != nil {
				panic(err)
			}
		}
	}
	

	http.Get(server_base_url + "/" + device_id + "/unregister")
}