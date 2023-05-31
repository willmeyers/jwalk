package jwalk

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Event struct {
	Id   string
	Name string
	Data interface{}
}

func OpenSSEventConnection(url string) (chan Event, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+"AUTH CODE")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	events := make(chan Event)

	go listen(reader, events)

	return events, nil
}

func listen(reader *bufio.Reader, events chan Event) {
	var event Event
	var buf bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		fmt.Println(string(line))
		if err != nil {
			close(events)
			return
		}

		switch {
		case bytes.HasPrefix(line, []byte(":")):
			// Comment, ignore
		case bytes.HasPrefix(line, []byte("retry:")):
			// TODO: Implement
		case bytes.HasPrefix(line, []byte("id: ")):
			event.Id = string(bytes.TrimSpace(line[4:]))
			fmt.Println("got", event.Id)
		case bytes.HasPrefix(line, []byte("event: ")):
			event.Name = string(line[7 : len(line)-1])
		case bytes.HasPrefix(line, []byte("data: ")):
			buf.Write(line[6:])
		case bytes.Equal(line, []byte("\n")):
			b := buf.Bytes()
			if bytes.HasPrefix(b, []byte("{")) {
				var data interface{}
				err := json.Unmarshal(b, &data)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
				event.Data = data
				buf.Reset()
				events <- event
				event = Event{}
			}
		default:
			fmt.Println("Unknown line: ", string(line))
		}
	}
}
