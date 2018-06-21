package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"flag"
	"errors"
	"github.com/fatih/color"
)

type State struct {
	Create 	bool
	Offset  bool
	Size	int
	Egg		string
}

func ParseCmdLine() *State {
	valid := true

	s := State{}

	flag.BoolVar(&s.Create, "c", false, "Create pattern")
	flag.BoolVar(&s.Offset, "o", false, "Find offset")
	flag.IntVar(&s.Size, "s", 20280, "Size of pattern")
	flag.StringVar(&s.Egg, "e", "", "Egg to hunt for")

	flag.Parse()

	if !s.Create && !s.Offset {
		fmt.Println("[!] You must use either the -c or -o flag")
		valid = false
	}
	if s.Create && s.Offset {
		fmt.Println("[!] You cannot use the -c and -o flag at the same time")
		valid = false
	}
	if s.Offset && s.Egg == "" {
		fmt.Println("[!] You must enter an egg to find using the -e flag ")
		valid = false
	}
	if valid {
		return &s
	}
	return nil
}

func Process(s *State){
	Ruler()
	if s.Create{
		pattern := CreatePattern(s)
		Ruler()
		fmt.Println(pattern)
	}

	if s.Offset{
		_, err := LocatePattern(s)
		if err != nil{
			color.Red(fmt.Sprint(err))
			Ruler()
		} else {
			Ruler()
		}
	}
}

func CreatePattern(s *State) string {
	charSet1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charSet2 := "abcdefghijklmnopqrstuvwxyz"
	charSet3 := "0123456789"

	fmt.Printf("[+] Creating pattern of %d bytes\n", s.Size)
	var pattern []string
	for len(pattern) < s.Size {
		for _, ch1 := range charSet1 {
			for _, ch2 := range charSet2 {
				for _, ch3 := range charSet3 {
					if len(pattern) < s.Size {
						pattern = append(pattern, string(ch1))
					}
					if len(pattern) < s.Size {
						pattern = append(pattern, string(ch2))
					}
					if len(pattern) < s.Size {
						pattern = append(pattern, string(ch3))
					}
				}
			}
		}
	}
	return strings.Join(pattern, "")
}

func LocatePattern(s *State) (int, error) {
	modes := []string{"normal", "upper", "lower"}
	extraText := ""
	comparePattern := ""

	for _, mode := range modes {
		if mode == "normal" {
			comparePattern = CreatePattern(s)
			extraText = " "
		} else if mode == "upper" {
			comparePattern = strings.ToUpper(CreatePattern(s))
			extraText = " (uppercase) "
		} else if mode == "lower" {
			comparePattern = strings.ToLower(CreatePattern(s))
			extraText = " (lowercase) "
		}

		if len(s.Egg) == 4 {
			asciiPat := s.Egg
			fmt.Printf("[+] Looking for egg %s in pattern of %d bytes\n", asciiPat, s.Size)

			if strings.Contains(comparePattern, asciiPat) {
				patIdx := strings.Index(comparePattern, asciiPat)
				output := fmt.Sprintf("[+] Egg pattern %s found in cyclic pattern%sat position %d\n", asciiPat, extraText, patIdx)
				color.Green(output)
				return patIdx, nil
			} else {
				//Reversed perhaps
				asciiPatRev := ReverseString(asciiPat)
				if strings.Contains(comparePattern, asciiPatRev) {
					patIdx := strings.Index(comparePattern, asciiPatRev)
					output := fmt.Sprintf("[+] Egg pattern %s (%s reversed) found in cyclic pattern%sat position %d\n", asciiPatRev, asciiPat, extraText, patIdx)
					color.Green(output)
					return patIdx, nil
				} else {
					return -1, errors.New(fmt.Sprintf("[!] Egg pattern %s not found in cyclic pattern%s\n", asciiPatRev, ""))
				}
			}
		}
		if len(s.Egg) == 10 {
			if s.Egg[:2] == "0x"{
				s.Egg = s.Egg[2:]
			} else {
				return -1, errors.New(fmt.Sprint("Invalid egg length"))
			}
		}
		if len(s.Egg) == 8 {
			hexPat := s.Egg
			bytePat, err := hex.DecodeString(hexPat)
			if err != nil {
				panic(err)
			}

			asciiPat := string(bytePat)
			fmt.Printf("[+] Looking for egg %s in pattern of %d bytes\n", asciiPat, s.Size)
			patIdx := strings.Index(comparePattern, asciiPat)
			if patIdx > -1 {
				output := fmt.Sprintf("[+] Egg pattern %s found in cyclic pattern%sat position %d\n", asciiPat, extraText, patIdx)
				color.Green(output)
				return patIdx, nil
			} else {
				//Reversed
				hexPatRev := ReverseHex(s.Egg)
				bytePatRev, err := hex.DecodeString(hexPatRev)
				if err != nil {
					panic(err)
				}

				asciiPatRev := string(bytePatRev)
				patIdx := strings.Index(comparePattern, asciiPatRev)
				if patIdx > -1 {
					output := fmt.Sprintf("[+] Egg pattern %s (%s reversed) found in cyclic pattern%sat position %d\n", asciiPatRev, asciiPat, extraText, patIdx)
					color.Green(output)
					return patIdx, nil
				} else {
					return -1, errors.New(fmt.Sprintf("[!] Egg pattern %s not found in cyclic pattern%s\n", asciiPat, extraText))
				}
			}
		}
	}
	return -1, errors.New(fmt.Sprint("[+] Egg not found\n"))
}

func ReverseHex(hexStr string) string {
	y := len(hexStr) - 1
	t := ""
	r := ""
	for y >= 0 {
		if y%2 != 0 {
			t += string(hexStr[y])
		} else {
			t = string(hexStr[y]) + t
			r += t
			t = ""
		}
		y--
	}
	return r
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Ruler() {
	fmt.Print("\n=====================================================\n")
}

func main() {
	state := ParseCmdLine()
	if state != nil {
		Process(state)
	}
}
