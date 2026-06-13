package gnuflag

import "fmt"

type ErrorInvalidOption struct {
	CommandName string
	OptionName  string
}

func (e *ErrorInvalidOption) Error() string {
	return fmt.Sprintf(`%s: invalid option -- "%s"
Try "%s --help" for more information.`, e.CommandName, e.OptionName, e.CommandName)
}

type ErrorNotRecognized struct {
	CommandName string
	OptionName  string
}

func (e *ErrorNotRecognized) Error() string {
	return fmt.Sprintf(`%s: unrecognized option "--%s"
Try "%s --help" for more information.`, e.CommandName, e.OptionName, e.CommandName)
}

type ErrorRequiresAnArg struct {
	CommandName string
	OptionName  string
}

func (e *ErrorRequiresAnArg) Error() string {
	return fmt.Sprintf(`%s: option requires an argument "--%s"
Try "%s --help" for more information.`, e.CommandName, e.OptionName, e.CommandName)
}

type ErrorRequired struct {
	CommandName string
	OptionName  string
}

func (e *ErrorRequired) Error() string {
	return fmt.Sprintf(`%s: missing required option "--%s"
Try "%s --help" for more information.`, e.CommandName, e.OptionName, e.CommandName)
}
