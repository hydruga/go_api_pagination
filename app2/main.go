package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Character struct {
	Info struct {
		Pages int     `json:"pages"`
		Next  *string `json:"next"`
		// We need pointers where the values could potentially be a nil state
	} `json:"info"`
	Results []struct {
		Name    string `json:"name"`
		Species string `json:"species"`
	} `json:"results"`
	Error string `json:"error"`
	// Add error type here for json error return
	// This is important as the http.Get request will not
	// return an error if the server responds with json error data
}

func main() {
	base_url := "https://rickandmortyapi.com/api/character"
	var totalAliens int

	character := Character{}
	character.GetJson(base_url)
	// // If we know the amount of pages we can use uncommented code
	//pages := character.Info.Pages
	var pages int
	//fmt.Println(pages)

	ch := make(chan int)
	// ch := make(chan int, pages)		// buffered channel
	done := make(chan struct{})

	go func() {
		read := 0

		// LOOP:
		//for i := range ch {
		for {
			read++
			// if read == pages {
			// 	break
			// }
			// if i == 0 {
			// 	goto LOOP
			// }
			character := Character{}
			//url := base_url + "?page=" + strconv.Itoa(i)
			url := base_url + "?page=" + strconv.Itoa(read)
			r, err := http.Get(url)
			fmt.Println("Retrieving ", url)
			if err != nil {
				//close the channel
				pages = read
				break
			}
			body, err := io.ReadAll(r.Body)

			json.Unmarshal(body, &character)
			if character.Error != "" {
				break
			}
			for _, item := range character.Results {
				if item.Species == "Alien" {
					totalAliens++
				}
			}
		}
		close(done)
	}()
	// for i := 0; i < pages; i++ {
	// 	ch <- i
	// }

	for i := 0; i < pages; i++ {
		ch <- i
	}
	<-done // block until ch is done
	close(ch)
	fmt.Println("Total aliens ", totalAliens)
}

func (c *Character) GetJson(url string) {
	r, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: Could not retrieve api data", err)
	}
	body, err := io.ReadAll(r.Body)

	json.Unmarshal(body, &c)

}
