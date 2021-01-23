package main

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/term"
)

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
