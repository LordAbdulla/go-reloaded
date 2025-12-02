package main

import (
    "fmt"
    "os"
    "reloaded"
    "strings"
)

func main() {
    args := os.Args[1:]
    if len(args) != 2 {
        os.Stdout.WriteString("Please enter a valid number of arguments. Usage: \"input.txt output.txt\"\n")
        return
    }

    text, err := os.ReadFile(args[0])
    if err != nil {
        os.Stdout.WriteString("Error reading file: " + args[0] + "\n")
        return
    }

    if !strings.HasSuffix(args[1], ".txt")  {
        panic("HEY, DONT MISS WITH THE FILES")
    }
    
    Content := string(text) // Here we have our text as a string

    Content = reloaded.Vowels(Content)
    Content = reloaded.Cases(Content)
    Content = reloaded.ULC(Content)
    Content = reloaded.Vowels(Content)
    Content = reloaded.Apostrophe(Content)
    Content = reloaded.Punctuation(Content)
    Content = reloaded.CleanText(Content)
    
    // Create a new file or truncate the existing file
    file, err := os.Create(args[1])
    if err != nil {
        fmt.Println("Error: creating file:", err)
        return
    }
    defer file.Close()

    // Write the string to the file
    _, err = file.WriteString(Content)
    if err != nil {
        fmt.Println("Error: writing to file:", err)
        return
    }

    //  Notice after successful write
    fmt.Println("File written successfully to", args[1])
}
