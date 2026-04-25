package main

import (
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ---------------- TOKENIZATION ----------------

func tokenize(text string) []string {
	var tokens []string
	var word strings.Builder

	flush := func() {
		if word.Len() > 0 {
			tokens = append(tokens, word.String())
			word.Reset()
		}
	}

	for i := 0; i < len(text); {
		r, size := utf8.DecodeRuneInString(text[i:])

		switch {

		case r == '(':
			flush()
			j := i
			for j < len(text) && text[j] != ')' {
				j++
			}
			if j < len(text) {
				tokens = append(tokens, text[i:j+1])
				i = j + 1
				continue
			}
			i += size
			word.WriteRune(r)
			continue

		case r == '\'':
			if word.Len() > 0 && i+size < len(text) {
				next, _ := utf8.DecodeRuneInString(text[i+size:])
				if unicode.IsLetter(next) || unicode.IsDigit(next) {
					word.WriteRune(r)
					i += size
					continue
				}
			}
			flush()
			tokens = append(tokens, "'")
			i += size
			continue

		case strings.ContainsRune(".,!?;:", r):
			flush()
			tokens = append(tokens, string(r))
			i += size
			continue

		case unicode.IsSpace(r):
			flush()
			i += size
			continue

		default:
			word.WriteRune(r)
			i += size
		}
	}

	flush()
	return tokens
}

// ---------------- MODIFIER ----------------

type modifier struct {
	cmd string
	n   int
}

var validCommands = map[string]bool{
	"up": true, "low": true, "cap": true,
	"hex": true, "bin": true,
}

func parseModifier(tok string) (modifier, bool) {
	if !strings.HasPrefix(tok, "(") || !strings.HasSuffix(tok, ")") {
		return modifier{}, false
	}

	inner := strings.TrimSpace(tok[1 : len(tok)-1])
	parts := strings.Split(inner, ",")

	cmd := strings.TrimSpace(parts[0])
	if !validCommands[cmd] {
		return modifier{}, false
	}

	n := 1
	if len(parts) > 1 {
		if v, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			n = v
		}
	}

	if n < 1 {
		n = 1
	}

	if cmd == "hex" || cmd == "bin" {
		n = 1
	}

	return modifier{cmd: cmd, n: n}, true
}

// ---------------- UTIL ----------------

func isWord(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func capitalize(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// ---------------- APPLY ----------------

func applyModifier(tokens []string, mod modifier) []string {
	if len(tokens) == 0 {
		return tokens
	}

	count := 0

	for i := len(tokens) - 1; i >= 0 && count < mod.n; i-- {
		if !isWord(tokens[i]) {
			continue
		}

		switch mod.cmd {

		case "up":
			tokens[i] = strings.ToUpper(tokens[i])

		case "low":
			tokens[i] = strings.ToLower(tokens[i])

		case "cap":
			tokens[i] = capitalize(tokens[i])

		case "hex":
			if val, err := strconv.ParseInt(tokens[i], 16, 64); err == nil {
				tokens[i] = strconv.Itoa(int(val))
			}

		case "bin":
			if val, err := strconv.ParseInt(tokens[i], 2, 64); err == nil {
				tokens[i] = strconv.Itoa(int(val))
			}
		}

		count++
	}

	return tokens
}

// ---------------- PROCESS ----------------

func process(tokens []string) []string {
	var result []string

	for _, tok := range tokens {
		if mod, ok := parseModifier(tok); ok {
			result = applyModifier(result, mod)
		} else {
			result = append(result, tok)
		}
	}

	return result
}

// ---------------- POST ----------------

func mergeQuotes(tokens []string) []string {
	var out []string

	for i := 0; i < len(tokens); i++ {
		if tokens[i] != "'" {
			out = append(out, tokens[i])
			continue
		}

		j := i + 1
		for j < len(tokens) && tokens[j] != "'" {
			j++
		}

		if j >= len(tokens) {
			out = append(out, tokens[i])
			continue
		}

		content := strings.Join(tokens[i+1:j], " ")
		merged := "'" + strings.TrimSpace(content) + "'"

		out = append(out, merged)
		i = j
	}

	return out
}

func mergePunctuation(tokens []string) []string {
	var out []string

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		if len(tok) == 0 {
			out = append(out, tok)
			continue
		}

		firstByte := tok[0]

		if strings.ContainsAny(string(firstByte), ".,!?;:") {
			group := tok
			j := i + 1
			for j < len(tokens) {
				next := tokens[j]
				if len(next) > 0 && strings.ContainsAny(string(next[0]), ".,!?;:") {
					group += next
					j++
				} else {
					break
				}
			}
			out = append(out, group)
			i = j - 1
		} else {
			out = append(out, tok)
		}
	}
	return out
}

func fixArticles(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ {
		word := tokens[i]
		next := tokens[i+1]

		if word != "a" && word != "A" {
			continue
		}

		clean := strings.Trim(next, ".,!?;:'\"()")
		if clean == "" {
			continue
		}

		firstByte := clean[0]

		if firstByte >= 'A' && firstByte <= 'Z' {
			firstByte += 32
		}

		if strings.ContainsAny(string(firstByte), "aeiouh") {
			if word == "A" {
				tokens[i] = "An"
			} else {
				tokens[i] = "an"
			}
		}
	}

	return tokens
}

func buildString(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(tokens[0])
	for i := 1; i < len(tokens); i++ {
		tok := tokens[i]
		if len(tok) > 0 && strings.ContainsAny(string(tok[0]), ".,!?;:") {
			b.WriteString(tok)
		} else {
			b.WriteByte(' ')
			b.WriteString(tok)
		}
	}
	return b.String()
}

func format(tokens []string) string {
	tokens = mergeQuotes(tokens)
	tokens = mergePunctuation(tokens)
	tokens = fixArticles(tokens)
	return buildString(tokens)
}

// ---------------- MAIN ----------------

func main() {
	if len(os.Args) != 3 {
		println("usage: go run . input.txt output.txt")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		println("read error:", err.Error())
		return
	}

	tokens := tokenize(string(data))
	processed := process(tokens)
	result := format(processed)

	err = os.WriteFile(os.Args[2], []byte(result), 0o644)
	if err != nil {
		println("write error:", err.Error())
		return
	}

	println("done")
}
