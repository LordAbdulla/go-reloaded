package reloaded

import (
    "regexp"
    "strings"
    "fmt"
)

// TrimSpaces removes spaces from both the beginning and the end of each line.
func TrimSpaces(text string) string {
    var result []string
    lines := strings.Split(text, "\n")
    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line) // Trim spaces at both ends
        result = append(result, trimmedLine)
    }
    return strings.Join(result, "\n")
}

// RemoveTrailingSpaces replaces multiple spaces or tabs with a single space.
func RemoveTrailingSpaces(text string) string {
    re := regexp.MustCompile(`[ \t]+`) // Regex to match multiple spaces or tabs
    text = re.ReplaceAllString(text, " ")
    text = TrimSpaces(text)
    return strings.TrimSpace(text)
}

// RemoveTrailingNewLines removes multiple newline characters and trims the result.
func RemoveTrailingNewLines(text string) string {
    re := regexp.MustCompile(`[\n]+`) // Regex to match multiple newlines
    text = re.ReplaceAllString(text, "\n")
    return strings.TrimSpace(text)
}

// CleanText performs a series of cleaning operations on the provided text.
func CleanText(text string) string {
    reSpaces := regexp.MustCompile(`  +`) // Regex to match two or more consecutive spaces
    if reSpaces.MatchString(text) {
        text = RemoveTrailingSpaces(text)
    }
    
    reNewlines := regexp.MustCompile(`\n\n+`) // Regex to match two or more consecutive newlines
    if reNewlines.MatchString(text) {
        text = RemoveTrailingNewLines(text)
    }
    
    reSpaceAtBeginOfNewline := regexp.MustCompile(`\n+ `) // Regex to match space at the beginning of a newline
    reSpaceAtbeginOfText := regexp.MustCompile(`\A +`)  // Regex to match spaces at the beginning of the text
    
    if len(text) > 1 && (reSpaceAtBeginOfNewline.MatchString(text) || reSpaceAtbeginOfText.MatchString(text)) {
        text = TrimSpaces(text)
    }

    if Noflag {
        fmt.Println(" Error: Invalid flags detected.")
    }
    return text
}
