package command

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

func (cmd Command) Definition() CommandDefinition {
	return cmd.commandDefinition
}

func (cmd Command) Parameters() []Parameter {
	return cmd.params
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

func (ctx *CommandContext) Command() Command {
	return ctx.command
}
