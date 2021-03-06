package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/user"
	"strings"
	"time"

	cursor "github.com/ahmetalpbalkan/go-cursor"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
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

	SetupCloseHandler()

	fmt.Print(cursor.Hide())

	homeDir = user.HomeDir

	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Print(cursor.ClearEntireScreen())

	startMenu()

	fmt.Print(cursor.Show())

}

func LoadText() {
	data, err := Asset("data/textdata")
	if err != nil {
		log.Fatal(err)
	}

	texts = strings.Split(string(data), "\n")

}

var uncorrected int

var characters, wordLength int
var elapsed time.Duration
var grosswpm, netwpm float64
var cpm int
var start time.Time

func startMenu() {
	for {
		choice := askChoice("Main Menu", "Offline (Practice)", "Multiplayer (WIP)", "Exit")
		if choice == 0 {
			LoadText()

			startGame()
		} else if choice == 1 {
			choice = askChoice("Multiplayer", "Join", "Host", "Back")
			switch choice {
			case 0:
				connect("localhost", port_default)
				break
			case 1:
				EnableKeyboard()
				port, name, err := ConfigureMultiplayer()
				if err != nil {
					log.Fatal(err)
				}
				DisableKeyboard()

				startServer(port, name)
				break
			case 2:
			default:
				continue
			}
		} else if choice == 2 {
			break
		}
	}
}

func startGame() {
	rightCharacters = 0

	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))

	text := texts[rand.Intn(len(texts))]
	split_text := strings.Split(text, " ")
	characters = len(text)
	wordLength = len(split_text)

	var textTyped string

	CustomPrint(text, "")
	fmt.Print(cursor.MoveDown(1))
	fmt.Print(cursor.MoveLeft(len(text)))

	DisableKeyboard()
	fmt.Print("\033[0mGame starting in \033[33m5")
	for i := 5; i > 0; i-- {
		cursor.MoveLeft(1)
		fmt.Print(cursor.MoveLeft(1))
		fmt.Print(i)
		time.Sleep(time.Second)
	}
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))
	EnableKeyboard()

	CustomPrint(text, "")

	start = time.Now()

	uncorrected = 0

	started := false

	for {
		ascii, _, err := GetChar()
		if err != nil {
			log.Fatal(err)
		}

		if started == false {
			started = true
			start = time.Now()
		}

		if ascii == 3 {
			break
		}

		if ascii == 127 {
			if len(textTyped)-1 >= 0 {
				textTyped = TrimLastChar(textTyped)
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
	fmt.Printf("\033[0mText (\033[94m%d\033[0m words) typed in \033[94m%s\033[0m.\n"+
		"Your Gross WPM is \033[94m%d\033[0m, Net WPM is \033[94m%d\033[0m and your CPM is \033[94m%d\033[0m.\n"+
		"You also typed with an accuracy of \033[94m%.2f%%\033[0m.\n"+
		"Press space to continue.\n",
		wordLength, elapsed.String(), int(grosswpm), int(netwpm), cpm, acc)

	for {
		char, _, _ := GetChar()
		if char == ' ' {
			break
		}
	}

	choice := askChoice("Do you wanna retry?", "Yes", "No")

	if choice == 0 {
		startGame()
	}
}

func updateTime() {
	elapsed = time.Since(start)
	grosswpm = float64(wordLength) / elapsed.Minutes()
	netwpm = grosswpm - (float64(uncorrected) / elapsed.Minutes())
	cpm = int(float64(characters) / elapsed.Minutes())
}

func CustomPrint(text, textWritten string) {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))

	term_width, _ := terminal.Width()
	textWr := wrap(text, int(term_width))

	uncorrected = 0
	for i, char := range textWritten {
		if i < len(textWritten) {
			if char == rune(text[i]) {
				fmt.Print("\033[32m" + string(textWr[i]) + "\033[0m")
			} else {
				uncorrected++
				fmt.Print("\033[41m" + string(textWr[i]) + "\033[0m")
			}
		}
	}

	for i := len(textWritten); i < len(text); i++ {
		if i == len(textWritten) {
			if string(rune(text[i])) != " " || string(rune(text[i])) != "\n" {
				fmt.Print("\033[45m" + string(rune(textWr[i])) + "\033[0m")
			} else {
				fmt.Print("\033[45m█\033[0m")
			}
		} else {
			fmt.Print(string(rune(textWr[i])))
		}
	}

	//fmt.Print(cursor.MoveUpperLeft(1) + cursor.MoveDown(1))
	//fmt.Printf("\033[1m\033[94mGross WPM: \033[97m%d  \033[94mNet WPM: \033[97m%d  \033[94mCPM: \033[97m%d\033[0m",
	//int(grosswpm), int(netwpm), cpm)
}
