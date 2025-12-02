# Text Formatter — Go Project

A Go-powered text editing tool that automatically corrects, formats, and transforms text based on inline instructions such as `(hex)`, `(bin)`, `(up)`, `(low)`, `(cap)`, punctuation correction, and grammar fixes.

---

##  Features

- **Hexadecimal to decimal conversion:**  
  `1E (hex)` → `30`

- **Binary to decimal conversion:**  
  `10 (bin)` → `2`

- **Text case transformations:**  
  - `(up)` → uppercase  
  - `(low)` → lowercase  
  - `(cap)` → capitalized  
  - Supports ranges: `(up, 3)`, `(low, 2)`, `(cap, 5)`

- **Punctuation fixing:**  
  Ensures correct spacing for: `. , ! ? : ;`  
  Keeps groups like `...` or `!?` intact.

- **Single quotes formatting:**  
  `' awesome '` → `'awesome'`  
  Supports multi-word quotes.

- **Grammar rule:**  
  Replace **a** with **an** if next word starts with a vowel or *h*.  
  Example: `a amazing` → `an amazing`

- **Works with input/output file arguments.**

---
## How to Clone

```bash
git clone https://github.com/LordAbdulla/go-reloaded
cd main/main.go .
