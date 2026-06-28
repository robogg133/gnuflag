package gnuflag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	CommandName  string
	receivedArgs argList

	Usage       string
	Description string
	Slogan      string
	info        []struct {
		Short []string
		Long  string
		Help  string
	}

	args []string

	lastid   uint8
	triggers map[string]struct {
		id     uint8
		isBool bool
	}

	parse map[uint8]func(string) error
}

type Trigger struct {
	Short, Full string
}

func NewParser(args []string) *Parser {
	return &Parser{
		CommandName:  args[0],
		receivedArgs: argList{src: args[1:]},
		triggers: make(map[string]struct {
			id     uint8
			isBool bool
		}),
		parse: make(map[uint8]func(string) error),
	}
}

func (p *Parser) addTrigger(full string, isBool bool, shorts ...string) uint8 {
	p.lastid++
	p.triggers[full] = struct {
		id     uint8
		isBool bool
	}{id: p.lastid, isBool: isBool}
	for _, short := range shorts {
		p.triggers[short] = struct {
			id     uint8
			isBool bool
		}{id: p.lastid, isBool: isBool}
	}
	return p.lastid
}

func (p *Parser) addHelp(full, help string, shorts ...string) {
	p.info = append(p.info, struct {
		Short []string
		Long  string
		Help  string
	}{Short: shorts, Long: full, Help: help})
}

func (p *Parser) NArgs() int {
	return len(p.args)
}

func (p *Parser) Arg(n int) string {
	return p.args[n]
}

func (p *Parser) SetFlagString(full, help string, value *string, shorts ...string) {
	validate(full, shorts...)
	id := p.addTrigger(full, false, shorts...)
	p.addHelp(full, help, shorts...)

	p.parse[id] = func(v string) error {
		*value = v
		return nil
	}
}

func (p *Parser) SetFlagInt(full, help string, value *int, shorts ...string) {
	validate(full, shorts...)
	id := p.addTrigger(full, false, shorts...)
	p.addHelp(full, help, shorts...)

	p.parse[id] = func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*value = n
		return nil
	}
}

func (p *Parser) SetFlagBool(full, help string, value *bool, shorts ...string) {
	validate(full, shorts...)
	id := p.addTrigger(full, true, shorts...)
	p.addHelp(full, help, shorts...)

	p.parse[id] = func(v string) error {
		if v == "false" {
			*value = false
		} else {
			*value = true
		}
		return nil
	}
}

func (p *Parser) Parse() error {
	p.help()
	for p.receivedArgs.Next() {
		arg := p.receivedArgs.Value()
		if s, ok := strings.CutPrefix(arg, "--"); ok {
			if len(s) == 0 {
				continue
			}
			value := ""
			if split := strings.SplitN(s, "=", 2); len(split) == 2 {
				arg = split[0]
				value = split[1]
			} else {
				arg = s
			}

			id, ok := p.triggers[arg]
			if !ok {
				err := &ErrorNotRecognized{CommandName: p.CommandName, OptionName: arg}
				fmt.Println(err.Error())
				os.Exit(2)
			}
			if id.isBool {
				err := p.parse[id.id]("")
				if err != nil {
					return err
				}
				continue
			}

			if value != "" {
				value = strings.TrimPrefix(value, `"`)
				value = strings.TrimSuffix(value, `"`)

				err := p.parse[id.id](value)
				if err != nil {
					return err
				}
			} else {
				if !p.receivedArgs.Next() {
					err := &ErrorRequiresAnArg{CommandName: p.CommandName, OptionName: arg}
					fmt.Println(err.Error())
					os.Exit(2)
				}

				err := p.parse[id.id](p.receivedArgs.Value())
				if err != nil {
					return err
				}
			}
		} else if s, ok := strings.CutPrefix(arg, "-"); ok {
			if len(s) == 0 {
				continue
			}
			arg = s
			id, ok := p.triggers[string(arg[0])]
			if !ok {
				err := &ErrorInvalidOption{CommandName: p.CommandName, OptionName: arg}
				fmt.Println(err.Error())
				os.Exit(2)
			}

			if id.isBool {
				err := p.parse[id.id](p.receivedArgs.Value())
				if err != nil {
					return err
				}
				continue
			}

			var value string
			if len(arg) > 1 {
				value = arg[1:]
				value = strings.TrimPrefix(value, `"`)
				value = strings.TrimSuffix(value, `"`)

			} else {
				if !p.receivedArgs.Next() {
					err := &ErrorRequiresAnArg{CommandName: p.CommandName, OptionName: arg}
					fmt.Println(err.Error())
					os.Exit(2)
				}
				value = p.receivedArgs.Value()
			}
			err := p.parse[id.id](value)
			if err != nil {
				return err
			}
		} else {
			p.args = append(p.args, arg)
			continue
		}
	}

	return nil
}

func (p *Parser) TryHelp() string {
	return fmt.Sprintf("Try \"%s --help\" for more information.", p.CommandName)
}

func validate(full string, shorts ...string) error {
	if len(shorts) == 0 && full == "" {
		return fmt.Errorf("invalid flags")
	}

	if strings.Contains(full, "=") {
		return fmt.Errorf("flag names cannot contain '='")
	}

	for _, short := range shorts {
		if strings.Contains(short, "=") {
			return fmt.Errorf("flag names cannot contain '='")
		}
	}

	return nil
}
