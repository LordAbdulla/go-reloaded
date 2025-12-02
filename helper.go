package reloaded

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Noflag bool

// This function handles flags with spaces, binary, hexadecimal, and other flags
func ULC(flag string) string {
	reULC := regexp.MustCompile(`(?i)(((\s|\n)\(\s*(cap|low|up|hex|bin)\s*\)))|((\s|\n)\(\s*(cap|low|up), \d+\s*\))`)

	for reULC.MatchString(flag) {
		pattern := regexp.MustCompile(`(?i)(?s)(.*?)\s+(\(\s*(cap|low|up|hex|bin)\s*\)|\(\s*(low|up|cap), (\d+)\s*\))`)
		match := pattern.FindString(flag)

		if len(match) < len(flag) && (flag[len(match)] != ' ' && flag[len(match)] != '\n' && flag[len(match)] != ')') {
			flag = ULC_flag(match) + " " + flag[len(match):] // to add a space in case I have a word sticking to the right of the flag
		} else {
			flag = ULC_flag(match) + flag[len(match):] // no need for space if I already have space
		}
	}
	return flag
}

func ULC_flag(s string) string {
	words := strings.Fields(s)
	flag := strings.ToLower(words[len(words)-1])

	s1 := s[:len(s)-len(flag)-1]

	// Validate hex and bin flags before any other processing
	if flag == "(hex)" && !IsHEx(FHex(words)) {
		fmt.Println(" Error: The provided hexadecimal value is invalid or out of range.")
		os.Exit(1)
	}

	if flag == "(bin)" && !IsBin(FBin(words)) {
		fmt.Println(" Error: The provided binary value is invalid or out of range.")
		os.Exit(1)
	}

	// Validate flags like (up, <non-numeric>) and (low, <non-numeric>) flags
	if flag == "(up," || flag == "(low," || flag == "(cap," {
		// Check if the flag number is valid
		if !IsValidInteger(words[len(words)-2]) {
			fmt.Printf(" Error: The number after the %s flag is not valid: %s\n", flag, words[len(words)-2])
			os.Exit(1)
		}
	}

	// Check if the flag is valid, otherwise remove it
	if ((flag == "(cap)" || flag == "(up)" || flag == "(low)") && FWord(words) == "") ||
		(flag == "(hex)" && FHex(words) == "") ||
		(flag == "(bin)" && FBin(words) == "") {
		Noflag = true
		return s1
	}

	// Processing valid flags
	switch flag {
	case "(cap)":
		return s1[:Index(s1, FWord(words))] + Capped(FWord(words)) + s1[Index(s1, FWord(words))+len(FWord(words)):]
	case "(up)":
		return s1[:Index(s1, FWord(words))] + strings.ToUpper(FWord(words)) + s1[Index(s1, FWord(words))+len(FWord(words)):]
	case "(low)":
		return s1[:Index(s1, FWord(words))] + strings.ToLower(FWord(words)) + s1[Index(s1, FWord(words))+len(FWord(words)):]
	case "(hex)":
		return s1[:Index(s1, FHex(words))] + Converter(FHex(words), 16, 10) + s1[Index(s1, FHex(words))+len(FHex(words)):]
	case "(bin)":
		return s1[:Index(s1, FBin(words))] + Converter(FBin(words), 2, 10) + s1[Index(s1, FBin(words))+len(FBin(words)):]
	}

	// handles flags with numbers
	temp := words[len(words)-1]                   // this is refers to the sec -->   2)
	flag2 := strings.ToLower(words[len(words)-2]) // len(words)-2 --> (2) is the index of the flag (cap,
	if temp[len(temp)-1] != ')' {
		return s
	}
	removeFlag := words[len(words)-2] + words[len(words)-1] // the flag is equal to: "(flag," + "num)"

	s2 := s[:len(s)-len(removeFlag)-2] // -2 because now we have two spaces to remove

	num, _ := strconv.Atoi(temp[:len(temp)-1]) // remove ")" from the number and convert it to int
	if len(Fflag(words, num)) < num {
		fmt.Printf(" Number Not Found, Applied the flag to the available words \n")
	}
	switch flag2 {
	case "(up,":
		for _, str := range Fflag(words, num) {
			s2 = UPLOW(s2, str, "(up,")
		}
		return s2
	case "(low,":
		for _, str := range Fflag(words, num) {
			s2 = UPLOW(s2, str, "(low,")
		}
		return s2
	case "(cap,":
		for _, str := range Fflag(words, num) {
			s2 = CAP(s2, str)
		}
		return s2
	}

	return s2
}

func Punctuation(text string) string {
	// Remove spaces before punctuation:
	re1 := regexp.MustCompile(` +([,.!?;:])`)
	text = re1.ReplaceAllString(text, "$1")

	// Ensure one space after punctuation: when what come after punctuation is not a punctuation or a whitespace.
	re2 := regexp.MustCompile(`([,.!?;:])([^,.'!?;:\s])`)
	text = re2.ReplaceAllString(text, "$1 $2")

	return text
}

func Vowels(text string) string {
	// Handle 'a' or 'A' followed by a word (inside quotes or backticks) starting with a vowel
	

	// Handle cases with 'a apple' or "a apple" or `a apple` (handles capital 'A' as well)
	reQuotes := regexp.MustCompile(`(?i)(\ba\b)\s*(["'` + "`" + `]?)\s*([aeiouh][a-zA-Z]*)`)
	text = reQuotes.ReplaceAllStringFunc(text, func(m string) string {
		parts := reQuotes.FindStringSubmatch(m)
		if len(parts) < 4 {
			return m
		}
		a := parts[1]
		quote := parts[2]
		word := parts[3]
		prefix := "an"
		if a == "A" {
			prefix = "AN"
		}
		if quote != "" {
			return prefix + " " + quote + word
		}
		return prefix + " " + word
	})

	// handle "a" or "A" followed by optional flags like (up,2) before a vowel-starting word
	reAWithFlags := regexp.MustCompile(`\b([Aa])\b((?:\s*$begin:math:text$\\s*(?:cap|low|up|hex|bin)(?:\\s*,\\s*[\\+\\-]?\\d+)?\\s*$end:math:text$)*)\s+([aeiouhAEIOUH][a-zA-Z]*)\b`)
	text = reAWithFlags.ReplaceAllStringFunc(text, func(m string) string {
		parts := reAWithFlags.FindStringSubmatch(m)
		if len(parts) < 4 {
			return m
		}
		a := parts[1]
		flags := parts[2] // can be empty
		word := parts[3]
		prefix := "an"
		if a == "A" {
			prefix = "AN"
		}
		// return "an" + flags + word without extra spaces
		return prefix + flags + " " + word
	})

	// Ensure the first letter "A" is turned into "AN" if followed by a vowel
	reCapA := regexp.MustCompile(`(?i)(\bA\b)\s?([aeiouh][a-zA-Z]*)\b`)
	text = reCapA.ReplaceAllString(text, "AN $2")

	return text
}

func Apostrophe(flag string) string {

	lines := strings.Split(flag, "\n")
	result := []string{}
	for _, line := range lines {
		re1 := regexp.MustCompile(`('\s+)`)   	 // Add a space before each apostrophe if not present
		line = re1.ReplaceAllString(line, " $1")
		re2 := regexp.MustCompile(`(\s+')`) 	// Add a space after each apostrophe if not present
		line = re2.ReplaceAllString(line, "$1 ")
		re3 := regexp.MustCompile(`\A'`)		// Add spaces around the apostrophe if it's at the beginning
		line = re3.ReplaceAllString(line, " ' ")	
		re4 := regexp.MustCompile(`'$`)				// Add spaces around the apostrophe if it's at the end of the line
		line = re4.ReplaceAllString(line, " ' ")
		
		// remove spaces on the right of even number apostrophes.
		re5 := regexp.MustCompile(`'\s+`) 	
		count := 0
		line = re5.ReplaceAllStringFunc(line, func(match string) string {
			if count%2 == 0 {
				count++
				return "'"
			} else {
				count++
				return match
			}
		})
		count = 0

		// remove spaces on the left of odd number apostrophes.
		re6 := regexp.MustCompile(`\s+'`)
		line = re6.ReplaceAllStringFunc(line, func(match string) string {
			if count%2 == 1 {
				count++
				return "'"
			} else {
				count++
				return match
			}
		})


		re7 := regexp.MustCompile(`[ ]+'`)		// Remove any spaces before the apostrophe and leave one space after it
		line = re7.ReplaceAllString(line, " '")
		re8 := regexp.MustCompile(`'[ ]+`)		// Remove any spaces after the apostrophe and leave one space before it
		line = re8.ReplaceAllString(line, "' ")
		re9 := regexp.MustCompile(`\A '`)		// If an apostrophe is at the start of the line, remove any leading spaces before it
		line = re9.ReplaceAllString(line, "'")

		result = append(result, strings.TrimRight(line, " \t"))
	}
	flag = strings.Join(result, "\n")

	return strings.TrimRight(flag, " \t")
}
