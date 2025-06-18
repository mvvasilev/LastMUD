package commandlib

type Command struct {
	commandDefinition CommandDefinition
	params            []Parameter
}

func CreateCommand(cmdDef CommandDefinition, parameters []Parameter) Command {
	return Command{
		commandDefinition: cmdDef,
		params:            parameters,
	}
}

func (cmd Command) Execute() (err error) {
	return cmd.commandDefinition.work(cmd.params...)
}

type commandContextError struct {
	err string
}

func createCommandContextError(err string) *commandContextError {
	return &commandContextError{
		err: err,
	}
}

func (cce *commandContextError) Error() string {
	return cce.err
}

type CommandContext struct {
	commandString string
	tokens        []Token

	command Command
}

func CreateCommandContext(commandRegistry *CommandRegistry, commandString string) (ctx *CommandContext, err error) {
	tokenizer := CreateTokenizer()

	tokens, tokenizerError := tokenizer.Tokenize(commandString)

	if tokenizerError != nil {
		err = tokenizerError
		return
	}

	commandDef := commandRegistry.Match(tokens)

	if commandDef == nil {
		err = createCommandContextError("Unknown command")
		return
	}

	params := commandDef.ParseParameters(tokens)

	ctx = &CommandContext{
		commandString: commandString,
		tokens:        tokens,
		command:       CreateCommand(*commandDef, params),
	}

	return
}

func (ctx *CommandContext) ExecuteCommand() (err error) {
	return ctx.command.Execute()
}
