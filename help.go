package gnuflag

import (
	"fmt"
	"os"
	"strings"
)

func (p *Parser) help() {
	id := p.addTrigger("help", true, "h")
	p.addHelp("help", "display this help and exit", "h")

	var builder strings.Builder
	builder.WriteString("Usage: ")
	builder.WriteString(p.CommandName + " ")
	builder.WriteString(p.Usage + "\n")
	builder.WriteString(p.Slogan + "\n\n")
	builder.WriteString(p.flagHelp())
	builder.WriteRune('\n')

	sep := separateString(p.Description, 75)
	for _, s := range sep {
		builder.WriteString(s)
		builder.WriteRune('\n')
	}
	helpmsg := builder.String()
	helpmsg = strings.TrimSuffix(helpmsg, "\n")
	p.parse[id] = func(_ string) error {
		fmt.Print(helpmsg)
		os.Exit(0)
		return nil
	}

}

func (p *Parser) flagHelp() string {
	var builder strings.Builder

	for _, info := range p.info {
		if info.Short == nil && info.Long != "" {
			builder.WriteString("      ")
			builder.WriteString("--" + info.Long)
			builder.WriteString("  ")
		} else {
			builder.WriteString("  ")
			lastI := len(info.Short) - 1
			totalSpaces := 0
			for i, short := range info.Short {
				n, err := builder.WriteString("-" + short)
				if err != nil {
					panic(err)
				}
				totalSpaces += n
				if i != lastI {
					n, err := builder.WriteRune(',')
					if err != nil {
						panic(err)
					}
					totalSpaces += n
				}
			}
			if info.Long != "" {
				n, err := builder.WriteString(", --" + info.Long)
				if err != nil {
					panic(err)
				}
				totalSpaces += n
			}
			for range 22 - totalSpaces {
				builder.WriteRune(' ')
			}
		}
		sep := separateString(info.Help, 44)
		builder.WriteString(sep[0])
		builder.WriteRune('\n')
		for _, s := range sep[1:] {
			builder.WriteString("                             ") // 30
			builder.WriteString(s)
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

func separateString(s string, width int) []string {
	if len(s) <= width {
		return []string{s}
	}

	result := make([]string, 0)
	bf := 0
	for i := 0; i < len(s); i += width {
		end := i + width
		end = min(end, len(s))

		afterSpace := offsetUntilNewLineOrSpace(s, end)
		result = append(result, s[bf:afterSpace])
		bf = afterSpace + 1
	}
	return result
}
func offsetUntilNewLineOrSpace(s string, offset int) int {
	for i := offset; i < len(s); i++ {
		if s[i] == '\n' || s[i] == ' ' {
			return i
		}
	}
	return len(s)
}
