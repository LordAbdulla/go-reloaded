package reloaded

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var NoFlag bool // Flag to track if no valid flag is found

func Cases(text string) string {

	if text == "" {
		fmt.Println("You provided an empty file.")
		return ""
	}

	// Check if there is no space after the comma for valid flags
	reNoSpaceAfterComma := regexp.MustCompile(`\(\s*(up|low|cap|hex|bin)\s*,([^\s\d\+\-\.\,])`) 
	if reNoSpaceAfterComma.MatchString(text) {
		fmt.Println("Error: Missing space after the comma or invalid character after the comma.")
		os.Exit(1) 
	}

	// Detect flags with invalid characters after them (non-numeric values like "rf-1")
	reInvalidFlagValues := regexp.MustCompile(`\(\s*(up|low|cap|hex|bin)\s*,\s*([^\d\+\-\s]+)\s*\)`)
	if reInvalidFlagValues.MatchString(text) {
		fmt.Println("Error: Invalid value after flag detected. Only numeric values are allowed.")
		os.Exit(1) 
	}

	// Detect flags surrounded by nested parentheses (multiple sets of parentheses)
	reNestedParentheses := regexp.MustCompile(`\(\(+\s*(up|low|cap|hex|bin)\s*\)+\)`) // Detect nested parentheses
	if reNestedParentheses.MatchString(text) {
		fmt.Println("Error: Nested parentheses detected around flags. Please remove extra parentheses.")
		os.Exit(1) 
	}

	reFlagHaveSpaces := regexp.MustCompile(`(?i)\(\s*(cap|low|up|hex|bin)\s*\)`) // This regex handles cases where there are spaces in the flags
	text = reFlagHaveSpaces.ReplaceAllString(text, "($1)")                      // It removes any spaces around the flag name, ensuring it’s in the format

	reFlagHaveSpaces = regexp.MustCompile(`\(\s*(low|up|cap),\s*(\d+)\s*\)`) // This regex handles flags with spaces followed by numbers
	text = reFlagHaveSpaces.ReplaceAllString(text, "($1, $2)")               // It removes extra spaces and ensures the format is

	reFlagToRight := regexp.MustCompile(`(?i)(\S)(\((cap|low|up|hex|bin)\)|\((low|up|cap),\s+(\d+)\))`) // This regex ensures there’s a space between a word and a flag when they’re next to each other (e.g., "word(cap)" becomes "word (cap)")
	text = reFlagToRight.ReplaceAllString(text, "$1 $2") // It adds the necessary space between the word and flag.

	reFlagSoloStart := regexp.MustCompile(`\A\s*(?i)(\((cap|low|up|hex|bin)\)|\((low|up|cap),\s+(\d+)\))`) // This regex checks if the flag appears at the beginning of the text without any word before it.
	if reFlagSoloStart.MatchString(text) {
		NoFlag = true                                      // Set NoFlag to true if a flag is at the beginning without a word before it.
		text = reFlagSoloStart.ReplaceAllString(text, "")  // It removes the flag from the start of the text.
	}

	reFlagWithSigns := regexp.MustCompile(`(?i)\(\s*(cap|low|up),\s*([\+\-]+)(\d+)\s*\)`) // This regex checks if the flag contains a negative or positive sign in the number (e.g., "(up, -5)" or "(low, +10)")
	text = reFlagWithSigns.ReplaceAllStringFunc(text, func(match string) string {
		// Check if the match contains a negative number
		reFlagNegativeNumber := regexp.MustCompile(`(?i)\((cap|low|up), ([\+\-]+)(\d+)\)`)
		if reFlagNegativeNumber.MatchString(match) {
			if hasNegativeSign(match) {
				fmt.Println("Error: Flag takes only positive numbers!")
				os.Exit(1)
			}
			return reFlagNegativeNumber.ReplaceAllString(match, "($1, $3)") 			// Remove negative sign and keep only the positive number
		}
		return match
	})

	reMultipleSpaces := regexp.MustCompile(`\(\s*(up|low|cap|hex|bin)\s*,\s*(\d+(\s+\d+)+)\s*\)`) // This regex detects flags with multiple spaces in the number (e.g., "(up, 9 9 9 99 9)")
	text = reMultipleSpaces.ReplaceAllStringFunc(text, func(match string) string {

		fmt.Println("Error: The flag contains multiple spaces in the number: ", match)
		os.Exit(1)
		return match
	})

	reEmptyFlag := regexp.MustCompile(`\(\s*(up|low|cap|hex|bin)\s*,\s*\)`) // This regex checks for empty flag values like "(hex, )" or "(bin, )"
	text = reEmptyFlag.ReplaceAllStringFunc(text, func(match string) string {
		fmt.Println("Error: The flag has an empty value: ", match)
		os.Exit(1)
		return match
	})

	reHexBin := regexp.MustCompile(`(?i)\(\s*(hex|bin)\s*,\s*([\+\-]\d+)\s*\)`) // This regex checks if a hex or bin flag has a negative number (e.g., "(bin, -1010)")
	text = reHexBin.ReplaceAllStringFunc(text, func(match string) string {
		fmt.Println("Error: The hex or bin flag cannot have negative numbers: ", match)
		os.Exit(1)
		return match
	})

	reLargeNumbers := regexp.MustCompile(`\(\s*(up|low|cap|hex|bin)\s*,\s*(\d+)\s*\)`) // This regex checks for large numbers in flags like "(up, 9438484784378347398483)"
	text = reLargeNumbers.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the number from the match
		reNumber := regexp.MustCompile(`\(\s*(up|low|cap|hex|bin)\s*,\s*(\d+)\s*\)`)
		matches := reNumber.FindStringSubmatch(match)

		if len(matches) > 2 {
			number := matches[2]                       // Get the number as a string
			_, err := strconv.ParseInt(number, 10, 64) // Try to parse it as a 64-bit integer

			if err != nil {
				fmt.Println("Error: Number out of range:", number)
				os.Exit(1)
			}
		}
		return match
	})

	return text
}

// Function to check if the flag contains a negative sign
func hasNegativeSign(match string) bool {
	re := regexp.MustCompile(`[\-]`) // This regex checks if there's a negative sign in the match.
	return re.MatchString(match)     
}
