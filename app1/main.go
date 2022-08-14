package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

func main() {
	url := "https://rickandmortyapi.com/api/character"
	c := Character{}
	c.GetJson(url)
	aliens := c.getTotalAliens()
	for {
		if c.Info.Next != nil {
			fmt.Println("Fetching next page ", *c.Info.Next)
			c.GetJson(*c.Info.Next)
			aliens += c.getTotalAliens()

		} else {
			break
		}
	}
	fmt.Println("Total aliens ", aliens)
}

func (character *Character) GetJson(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: Issues fetching url ", err)

	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Cannot read response body ", err)

	}
	if err := json.Unmarshal([]byte(body), &character); err != nil {
		fmt.Println("Error: Cannot unmarshal JSON ", err)
	}

}

func (character *Character) getTotalAliens() int {
	aliens := 0
	for _, val := range character.Results {
		if val.Species == "Alien" {
			aliens++
		}
	}
	return aliens

}
