package reloaded

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Rune Finder
func Index(s, f string) int {
	Search := []rune(s)
	Find := []rune(f)
	for i := len(Search) - len(Find); i >= 0; i-- {
		if string(Search[i:i+len(Find)]) == f {
			return i 
		}
	}
	return -1 // Not found
}

// Letter Finder
func FLetter(s string) bool {
	for _, char := range s {
		if unicode.IsLetter(char) {
			return true
		}
	}
	return false
}

// Word Finer in a string
func FWord(words []string) string {
	idx := len(words) - 1
	word := ""
	for i := idx - 1; i >= 0; i-- {
		if FLetter(words[i]) {			
			word = words[i]
			break
		}
	}
	return word
}

func IsValidInteger(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

// Hex Finder in a string 
func FHex(words []string) string {
	idx := len(words) - 1
	for i := idx - 1; i >= 0; i-- {
		if IsHEx(words[i]) {			
			return words[i]
		}
	}
	return ""
}

// Is Valid Hex
func IsHEx(str string) bool {
	re := regexp.MustCompile(`[^0-9A-Fa-f]`)
	if re.MatchString(str) {
		return false
	}
	
	var num int64
	n, err := fmt.Sscanf(str, "%x", &num)
	if err != nil || n != 1 {
		return false
	}

	// Maximum value for a 64-bit signed integer
	maxValue := int64(9223372036854775807)
	if num > maxValue {
		return false
	}

	return true
}

// Bin Finder
func FBin(words []string) string {
	idx := len(words) - 1
	for i := idx; i >= 0; i-- {
		if IsBin(words[i]) {
			return words[i]
		}
	}
	return ""
}

// Is Valid bin
func IsBin(str string) bool {
	re := regexp.MustCompile(`[^01]`)
	if re.MatchString(str) {
		return false
	}

	var num int64
	n, err := fmt.Sscanf(str, "%b", &num)
	if err != nil || n != 1 {
		return false
	}

	// Maximum value for a 64-bit signed integer
	maxValue := int64(9223372036854775807)
	if num > maxValue {
		return false
	}

	return true
}

// find the word before the flag
func Fflag(words []string, n int) []string {
	idx := len(words) - 3 		//this removes the flag and the num
	flagsf := []string{}
	if n <= 0 {
		fmt.Println("Error: No word can be found before the flag")
		return flagsf
	}
	for i := idx; i >= 0; i-- {
		if FLetter(words[i]) {
			flagsf = append(flagsf, words[i])
			n--
			if n <= 0 {
				break
			}
		}
	}
	if len(flagsf) == 0 {
		Noflag = true
	}
	return flagsf
}

// find the flag and replace it (up, low)
func UPLOW(s, f, flag string) string {
	Search := []rune(s)
	Find := []rune(f)
	var replace []rune
	switch flag {
	case "(up,":
		replace = []rune(strings.ToUpper(string(Find)))
	case "(low,":
		replace = []rune(strings.ToLower(string(Find)))
	}

	for i := len(Search) - len(Find); i >= 0; i-- {
		if string(Search[i:i+len(Find)]) == f {
			for j := 0; j < len(Find); j++ {
				Search[i+j] = replace[j]

			}
			return string(Search)
		}
	}
	return s 
}

// find the flag and replace it (cap)
func CAP(s, f string) string {
	Search := []rune(s)
	Find := []rune(f)
	replacement := []rune(Capped(strings.ToLower(f)))

	for i := len(Search) - len(Find); i >= 0; i-- {
		if string(Search[i:i+len(Find)]) == f && (i == 0 || !unicode.IsLetter(Search[i-1])) && (i+len(Find) == len(Search) || !unicode.IsLetter(Search[i+len(Find)])) {
			for j := 0; j < len(Find); j++ {
				Search[i+j] = replacement[j]
			}
			return string(Search)
		}
	}
	return s 
}

// Converter --> converts a base to base 
func Converter(s string, a, b int) string {
	num, err := strconv.ParseInt(s, a, 64)
	if err != nil {
		return s 
	}
	return strconv.FormatInt(num, b)
}

// Capping
func Capped(s string) string {
	s = strings.ToLower(s)
	Search := []rune(s)
	for i, char := range Search {
		if unicode.IsLetter(char) {
			Search[i] = unicode.ToUpper(char) 
			break
		}
	}
	return string(Search)
}

// Checks if the number of the flag is negative or positive 
func PosNegFlag(flag string, reFlagNegativeNumber *regexp.Regexp) bool {
	matches := reFlagNegativeNumber.FindStringSubmatch(flag)
	signs := matches[2] 
	PosCount := 0
	NegCount := 0

	for _, char := range signs {
		switch char {
		case '+':
			PosCount++
		case '-':
			NegCount++
		}
	}
	// Determine the sign based on the counts of + and -.
	// If the number of - signs is odd, the result is negative; otherwise, it's positive.
	return NegCount%2 == 0
}
