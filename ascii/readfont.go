package asciiart

import (
	"bufio"
	"os"
)

//array of all kind of font name. you can expend it
var validFonts = [3]string{"standard", "shadow", "thinkertoy"}

/*
Opens the file the selected fonts characters per 8 lines, separating them by newline and sorts it into an 2d array.
*/

func ReadFontFile(fontname string) ([95][8]string, error) {

	var fontarray = [95][8]string{}

	font := "ascii/fonts/" + fontname + ".txt"
	fontFile, err := os.Open(font)
	if err != nil {
		return [95][8]string{}, err
	}

	scanner := bufio.NewScanner(fontFile)

	// startline is empty so need to start at -1

	currentchar := -1
	currentline := 0

	// iterate over lines in the fontfile
	for scanner.Scan() {
		// if line is empty its the start of a new char
		if scanner.Text() == "" {
			currentchar++
			currentline = 0
		} else {
			fontarray[currentchar][currentline] = scanner.Text() // fill the array
			currentline++
		}

	}

	if err := scanner.Err(); err != nil {
		fontFile.Close()
		return [95][8]string{}, err
	}

	err = fontFile.Close() // proper close instead of defer which apparently causes bugs

	return fontarray, err

}

func FontValidation(fontname string) bool {
	for _, font := range validFonts {
		if fontname == font {
			return true
		}
	}
	return false
}
