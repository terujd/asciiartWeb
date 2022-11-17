package asciiart

import (
	"errors"
	"fmt"
	"strings"
)

type RequestError struct {
	StatusCode int

	Err error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: %v", r.StatusCode, r.Err)
}
func doBadRequest() error {
	return &RequestError{
		StatusCode: 400,
		Err:        errors.New("wrong input"),
	}
}
func doMissingRequest() error {
	return &RequestError{
		StatusCode: 404,
		Err:        errors.New("missing file"),
	}

}

func Printer(input string, fontname string) (string, error) {

	fontArray, err := ReadFontFile(fontname)
	if err != nil {
		return "", doMissingRequest()
	}
	input = strings.ReplaceAll(input, `\n`, "\n")

	chararray := make([]string, 0)
	for _, c := range input {
		if c >= 32 && c <= 126 {
			chararray = make([]string, 8)
			break
		}
	}

	for i := 0; i < len(input); i++ {
		// Get the valid character from the font
		if input[i] != '\n' && input[i] >= 32 && input[i] <= 126 {
			runeChar := rune(input[i])
			// Need to substract 32 to get the right position in the ascii table
			char := int(runeChar) - 32
			output := fontArray[char]

			for linePos, line := range output {
				chararray[linePos+len(chararray)-8] += line
			}
		} else if input[i] == '\n' || input[i] == 13 {
			// If there is no character after the newline, add 1 line, otherwise add 8

			if i == 0 || i == len(input)-1 || input[i+1] == '\n' {
				chararray = append(chararray, make([]string, 1)...)
			} else {
				chararray = append(chararray, make([]string, 8)...)
			}
		} else {

			fmt.Println("Error: Invalid character")
			return "", doBadRequest()
		}
	}
	result := ""

	result = strings.Join(chararray, "\n")
	return result, nil
}
