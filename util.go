package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/pkg/term"
)

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\033[0m")
		fmt.Println(cursor.Show())
		EnableKeyboard()
		log.Println("^C detected. Force closing...")
		os.Exit(0)
	}()
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

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if size < 1 {
		return s
	}
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}

	return s[:len(s)-size]
}

func GetChar() (ascii int, keyCode int, err error) {
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

func askChoice(prompt string, choices ...string) int {
	choice := 0

	for {
		fmt.Print(cursor.ClearEntireScreen(),
			cursor.MoveUpperLeft(1))

		if prompt != "" {
			fmt.Println(prompt)
		}

		for i := 0; i < len(choices); i++ {
			if i == choice {
				fmt.Print("\033[94m>\033[0m ")
			} else {
				fmt.Print("  ")
			}

			fmt.Println(choices[i])
		}
		ascii, keyCode, err := GetChar()
		if err != nil {
			log.Fatal(err)
		}

		if rune(ascii) == 'j' || keyCode == 40 {
			if choice+1 < len(choices) {
				choice++
			} else {
				choice = 0
			}
		}
		if rune(ascii) == 'k' || keyCode == 38 {
			if choice-1 > -1 {
				choice--
			} else {
				choice = len(choices) - 1
			}
		}

		if rune(ascii) == 'l' || ascii == 13 || keyCode == 39 {
			break
		}
	}

	fmt.Print(cursor.ClearEntireScreen(),
		cursor.MoveUpperLeft(1))
	return choice
}

func PrintMessage(message string) {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveUpperLeft(1))
	fmt.Println(message)
	fmt.Println("Press any key to continue.")
	GetChar()
}

func DisableKeyboard() {
	c := exec.Command("stty -echo")
	c.Stdout = os.Stdout
	c.Run()
}

func EnableKeyboard() {
	c := exec.Command("stty echo")
	c.Stdout = os.Stdout
	c.Run()
}
