// Server running at litbit.in

package main

import (
	"fmt"
	"net/http"
	"regexp"
	"html/template"
)

type Device struct {
	Id string
	Queue []string
}

var list_regex = regexp.MustCompile(`\/list\/?$`)
var app_regex = regexp.MustCompile(`\/([A-Z]{4})\/?$`)
var register_regex = regexp.MustCompile(`\/([A-Z]{4})\/register\/?$`)
var unregister_regex = regexp.MustCompile(`\/([A-Z]{4})\/unregister\/?$`)
var get_regex = regexp.MustCompile(`\/([A-Z]{4})\/get\/?$`)
var add_regex = regexp.MustCompile(`\/([A-Z]{4})\/add\/?$`)

var devices []Device

func remove(s []string, i int) []string {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func remove_device(s []Device, i int) []Device {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func matchUrl(w http.ResponseWriter, r *http.Request) {
	switch {
	case list_regex.MatchString(r.URL.Path):
		for _, dev := range devices {
			fmt.Fprintf(w, "%s\n", dev.Id)
			for _, url := range (dev.Queue) {
				fmt.Fprintf(w, "%s\n", url)
			}
		}
		return

	case unregister_regex.MatchString(r.URL.Path):
		match := unregister_regex.FindStringSubmatch(r.URL.Path)

		for i := 0; i < len(devices); i++ {
			if devices[i].Id == match[1] {
				devices = remove_device(devices, i)

				fmt.Fprintf(w, "Device id %s has been unregistered\n", match[1])
				fmt.Printf("Device id %s has been unregistered\n", match[1])
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request\n")
		return

	case register_regex.MatchString(r.URL.Path):
		match := register_regex.FindStringSubmatch(r.URL.Path)

		for _, dev := range devices {
			if dev.Id == match[1] {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "403 Forbidden\n")
				return
			}
		}

		new_device := Device {
			Id: match[1],
		}

		devices = append(devices, new_device)

		fmt.Printf("Registered device with id %s\n", match[1])
		fmt.Fprintf(w, "Registered device with id %s\n", match[1])
		return

	case get_regex.MatchString(r.URL.Path):
		match := get_regex.FindStringSubmatch(r.URL.Path)

		for i := 0; i < len(devices); i++ {
			if devices[i].Id == match[1] {
				//fmt.Printf("Device id %s requested next song\n", match[1])
				if len(devices[i].Queue) > 0 {
					fmt.Fprintf(w, "%s\n", devices[i].Queue[0])
					devices[i].Queue = remove(devices[i].Queue, 0)
					return
				} else {
					w.WriteHeader(http.StatusNoContent)
					fmt.Fprintf(w, "204 No Content\n")
					return
				}
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request\n")
		return

	case add_regex.MatchString(r.URL.Path):
		match := add_regex.FindStringSubmatch(r.URL.Path)

		if r.Method == "POST" {
			for i := 0; i < len(devices); i++ {
				if devices[i].Id == match[1] {
					r.ParseForm()
					devices[i].Queue = append(devices[i].Queue, r.PostFormValue("url"))
					fmt.Printf("Added song to device id %s\n", match[1])

					http.ServeFile(w, r, "www/app-add.html")
					return
				}
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request\n")
		return
			
	case app_regex.MatchString(r.URL.Path):
		match := app_regex.FindStringSubmatch(r.URL.Path)

		for _, dev := range devices {
			if dev.Id == match[1] {
				template_vars := struct {
					Id string
				} {
					Id: dev.Id,
				}

				t, err := template.ParseFiles("www/app.html")
				if err != nil {
					panic(err)
				}
				err = t.Execute(w, template_vars)
				if err != nil {
					panic(err)
				}
			}
		}
		return
			
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request\n")
		return
	}	
}

func main() {
	fmt.Printf("Listening on :80\n")

	http.HandleFunc("/", matchUrl)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
} 