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

func (character *Character) GetJson(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("Error encountered %s :", err)
		fmt.Println(err)
		return err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Cannot read response body ", err)
		return err
	}
	if err := json.Unmarshal([]byte(body), &character); err != nil {
		fmt.Println("Error: Cannot unmarshal JSON ", err)
		return err
	}
	return nil
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
