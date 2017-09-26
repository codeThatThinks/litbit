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
	"github.com/xdsopl/framebuffer/src/framebuffer"
	"image"
	"image/draw"
	"image/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype"
)

var id_length = 4
var server_base_url = "http://litbit.in"
var device_id string
var fb draw.Image
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

func DrawMessage() {
	ClearScreen()

	context := freetype.NewContext()
	context.SetDst(fb)
	context.SetSrc(&image.Uniform{color.RGBA{255, 255, 255, 255}})
	context.SetDPI(170)
	context.SetHinting(font.HintingFull)
	context.SetClip(fb.Bounds())
	goregular, _ := freetype.ParseFont(goregular.TTF)
	gomedium, _ := freetype.ParseFont(gomedium.TTF)

	context.SetFontSize(11)
	context.SetFont(goregular)
	context.DrawString("Go to", fixed.P(fb.Bounds().Min.X + 125, fb.Bounds().Min.Y + 80))
	context.DrawString("to keep the party goin", fixed.P(fb.Bounds().Min.X + 40, fb.Bounds().Min.Y + 170))

	context.SetFontSize(13)
	context.SetFont(gomedium)
	context.DrawString(server_base_url + "/" + device_id, fixed.P(fb.Bounds().Min.X + 20, fb.Bounds().Min.Y + 125))
}

func ClearScreen() {
	draw.Draw(fb, fb.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
}

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

	device_id = ""

	for !registered {
		// gen unique id
		device_id = GenerateDeviceId()

		// try to register
		var resp *http.Response
		var err error
		resp, err = http.Get(server_base_url + "/" + device_id + "/register")
		for err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Retrying...\n")
			resp, err = http.Get(server_base_url + "/" + device_id + "/register")
		}

		if resp.StatusCode == 200 {
			registered = true
		} else if resp.StatusCode != 403 {
			fmt.Printf("Error: Recieved status code %i when trying to register", resp.StatusCode)
			os.Exit(1)
		}
	}

	fb, _ = framebuffer.Open("/dev/fb1")
	handleExit()

	DrawMessage()
	fmt.Printf("Go to %s/%s to request songs\n", server_base_url, device_id)

	for !shouldQuit {
		var resp *http.Response
		var err error
		resp, err = http.Get(server_base_url + "/" + device_id + "/get")
		for err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Retrying...\n")
			resp, err = http.Get(server_base_url + "/" + device_id + "/get")
		}

		if resp.StatusCode == 200 {
			// play the song
			urlb, _ := ioutil.ReadAll(resp.Body)
			url := strings.TrimSpace(string(urlb))
			fmt.Printf("Playing %s\n", url)
			command := exec.Command("/bin/sh", "vlc.sh", url)
			_ = command.Run()
		}
	}

	http.Get(server_base_url + "/" + device_id + "/unregister")
	fmt.Printf("\n");
	ClearScreen()
}
