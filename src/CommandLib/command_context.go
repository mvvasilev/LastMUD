package commandlib

type commandContext struct {
	commandString string
	tokens        []Token

	command Command
}

func CreateCommandContext(commandString string) (ctx *commandContext, err error) {
	tokenizer := CreateTokenizer()

	tokens, tokenizerError := tokenizer.Tokenize(commandString)

	if tokenizerError != nil {
		err = tokenizerError
		return
	}

	ctx = &commandContext{
		commandString: commandString,
		tokens:        tokens,
	}

	return
}

func (ctx *commandContext) Execute() (err error) {
	ctx.command.Execute()
}
