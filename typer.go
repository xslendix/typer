package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os/user"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/pkg/term"
)

var texts []string
var rightCharacters int

var homeDir string

var timer time.Timer

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	homeDir = user.HomeDir

	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Print(cursor.ClearEntireScreen())

	LoadText()

	startGame()

}

func LoadText() {
	content, err := ioutil.ReadFile(homeDir + "/.local/share/textdata")
	if err != nil {
		fmt.Println(err)
	}

	texts = strings.Split(string(content), "\n")

}

var uncorrected int

var characters, wordLength int
var elapsed time.Duration
var grosswpm, netwpm float64
var cpm int
var start time.Time

func startGame() {
	rightCharacters = 0

	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))

	text := texts[rand.Intn(len(texts))]
	split_text := strings.Split(text, " ")
	characters = len(text)
	wordLength = len(split_text)

	var textTyped string

	fmt.Print("\033[0mGame starting in \033[33m5")
	for i := 5; i > 0; i-- {
		cursor.MoveLeft(1)
		fmt.Print(cursor.MoveLeft(1))
		fmt.Print(i)
		time.Sleep(time.Second)
	}
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))

	CustomPrint(text, "")

	start = time.Now()

	uncorrected = 0

	for {
		ascii, _, err := getChar()
		if err != nil {
			log.Fatal(err)
		}

		if ascii == 3 {
			break
		}

		if ascii == 127 {
			if len(textTyped)-1 >= 0 {
				textTyped = trimLastChar(textTyped)
			}
		} else if ascii == 8 {
			if len(textTyped)-1 >= 0 {
				RemoveLastWord(&textTyped)
			}
		} else {
			if len(textTyped) != len(text) {
				textTyped = textTyped + string(rune(ascii))
			}
		}

		if len(textTyped) != len(text) && len(textTyped)-1 >= 0 {
			if textTyped[len(textTyped)-1] == text[len(textTyped)-1] {
				rightCharacters++
			}
		}

		//updateTime()

		CustomPrint(text, textTyped)

		if len(text) == len(textTyped) {
			break
		}

	}

	updateTime()

	acc := (netwpm / grosswpm) * 100
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))
	fmt.Printf("\033[0mText typed in \033[94m%s\033[0m. Your Gross WPM is \033[94m%d\033[0m, Net WPM is \033[94m%d\033[0m and your CPM is \033[94m%d\033[0m.\n"+
		"You also typed with an accuracy of \033[94m%.2f%%\033[0m.\n", elapsed.String(), int(grosswpm), int(netwpm), cpm, acc)

}

func updateTime() {
	elapsed = time.Since(start)
	grosswpm = float64(wordLength) / elapsed.Minutes()
	netwpm = grosswpm - (float64(uncorrected) / elapsed.Minutes())
	cpm = int(float64(characters) / elapsed.Minutes())
}

func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 2 && bytes[0] == 17 {
		ascii = int(bytes[1])
		keyCode = 17
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}

func CustomPrint(text, textWritten string) {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))

	uncorrected = 0
	for i, char := range textWritten {
		if i < len(textWritten) {
			if char == rune(text[i]) {
				fmt.Print("\033[32m" + string(char) + "\033[0m")
			} else {
				uncorrected++
				fmt.Print("\033[41m" + string(text[i]) + "\033[0m")
			}
		}
	}

	for i := len(textWritten); i < len(text); i++ {
		if i == len(textWritten) {
			if string(rune(text[i])) != " " {
				fmt.Print("\033[45m" + string(rune(text[i])) + "\033[0m")
			} else {
				fmt.Print("\033[45m█\033[0m")
			}
		} else {
			fmt.Print(string(rune(text[i])))
		}
	}

	//fmt.Print(cursor.MoveUpperLeft(1) + cursor.MoveDown(1))
	//fmt.Printf("\033[1m\033[94mGross WPM: \033[97m%d  \033[94mNet WPM: \033[97m%d  \033[94mCPM: \033[97m%d\033[0m",
	//int(grosswpm), int(netwpm), cpm)
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if size < 1 {
		return s
	}
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}

	return s[:len(s)-size]
}

func RemoveLastWord(s *string) {
	words := strings.Split(*s, " ")
	if len(words) > 0 {
		words = words[:len(words)-1]
	}
	*s = strings.Join(words, " ")
	if len(*s) != 0 {
		*s += " "
	}
}
