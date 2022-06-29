package remotedata

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func readBody(resp http.Response) ([]byte, error) {
	// Create a buffer, read 32 bytes at a time
	byte_buffer := make([]byte, 32)
	file_contents := make([]byte, 0)
	for {
		bytes_read, err := resp.Body.Read(byte_buffer)
		if bytes_read > 0 {
			file_contents = append(file_contents, byte_buffer[:bytes_read]...)
		}
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to read body:\n%s", err)
			} else {
				break
			}
		}
	}
	return file_contents, nil
}

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, arr := readBody(*resp)
	if arr != nil {
		return nil, arr
	}
	return data, nil
}
