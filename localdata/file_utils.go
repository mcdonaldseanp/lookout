package localdata

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/clibuild/validator"
)

const STDIN_IDENTIFIER string = "__STDIN__"

func readFromStdin() string {
	var builder strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}
	return builder.String()
}

func ChooseFileOrStdin(specfile string, use_stdin bool) (string, error) {
	if use_stdin {
		if len(specfile) > 0 {
			return "", &errtype.InvalidInput{
				Message: "cannot specify both a file and to use stdin",
				Origin:  nil,
			}
		}
		return STDIN_IDENTIFIER, nil
	} else {
		// Validate that the thing is actually a file on disk before
		// going any further
		//
		// Cheat a little with the validator: this function is mostly used
		// for the CLI commands, so use a name that shows it's the flag
		err := validator.ValidateParams(fmt.Sprintf(
			`[{"name":"--file","value":"%s","validate":["NotEmpty","IsFile"]}]`,
			specfile,
		))
		if err != nil {
			return "", err
		}
		return specfile, nil
	}
}

func ReadFileOrStdin(maybe_file string) ([]byte, error) {
	var raw_data []byte
	if maybe_file == STDIN_IDENTIFIER {
		raw_data = []byte(readFromStdin())
	} else {
		// raw_data was already created so you have to define err now too
		var err error
		raw_data, err = ReadFileInChunks(maybe_file)
		if err != nil {
			return nil, err
		}
	}
	return raw_data, nil
}

func ReadFileInChunks(location string) ([]byte, error) {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open file:\n%s", err)
	}
	defer f.Close()

	// Create a buffer, read 32 bytes at a time
	byte_buffer := make([]byte, 32)
	file_contents := make([]byte, 0)
	for {
		bytes_read, err := f.Read(byte_buffer)
		if bytes_read > 0 {
			file_contents = append(file_contents, byte_buffer[:bytes_read]...)
		}
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to read file:\n%s", err)
			} else {
				break
			}
		}
	}
	return file_contents, nil
}

func OverwriteFile(location string, data []byte) error {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file:\n%s", err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file:\n%s", err)
	}
	return nil
}
