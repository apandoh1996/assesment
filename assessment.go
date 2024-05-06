package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	data, err := os.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	// added Wait group to finish all go routines before exiting
	var wg sync.WaitGroup
	for _, line := range lines {
		// it represents number of go routines before the Waitgroup becomes unblocked
		wg.Add(1)
		// passed line as an argument to avoid race condition
		go func(line string) {
			defer wg.Done()
			data, err := getData(line)
			if err != nil {
				log.Printf("unable to get status: %v", err)
			}
			fmt.Println(data)
			if data.(string) == "foo" {
				fmt.Printf("data found: %s", data)
				os.Exit(0)
			}
		}(line)
	}
	wg.Wait()
}

func getData(line string) (interface{}, error) {
	// the string's proper structure according to input.txt is used in order to add exact URL inside the URL variable
	var location struct {
		URL string `json:"location"`
	}
	if err := json.Unmarshal([]byte(line), &location); err != nil {
		return "", err
	}
	fmt.Printf(location.URL)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, location.URL, nil)
	if err != nil {
		return "", err
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, res.Body); err != nil {
		return "", err
	}

	var payload struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		return "", err
	}

	return payload.Data, nil
}
