package world

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"github.com/google/uuid"
)

// Rules for output:
// - All server output is prepended with "< "
// - Server output is max 80 characters per line ( including th above prefix )
// - If longer than 80 characters, break up into multiple lines
// - Following lines of same output are not prepended with "< ", but have an offset to account for it
// - Bottom-most line of telnet must always start with "> ", and if input currently exists there, must be preserved in case of server output

//const (
//	maxLineLength = 80
//	prefix        = "< "
//	continuation  = "  " // 2-space offset to align with prefix
//)

//func formatServerOutput(text string) []string {
//	// Strip newline characters to avoid confusion
//	text = strings.ReplaceAll(text, "\n", " ")
//	words := strings.Fields(text)
//
//	var lines []string
//	line := prefix
//
//	for _, word := range words {
//		// Check if word fits on current line
//		if len(line)+len(word)+1 > maxLineLength {
//			lines = append(lines, line)
//			// Begin new line with continuation indent
//			line = continuation + word
//		} else {
//			if len(line) > len(continuation) {
//				line += " " + word
//			} else {
//				line += word
//			}
//		}
//	}
//
//	if len(line) > 0 {
//		lines = append(lines, line)
//	}
//
//	return lines
//}
//
//func formatServerOutputPreservingInput(output string) []byte {
//	lines := formatServerOutput(output)
//
//	saveCursor := "\x1b7"
//	restoreCursor := "\x1b8"
//	moveCursorUp := func(n int) string {
//		return fmt.Sprintf("\x1b[%dA", n)
//	}
//	eraseLine := "\x1b[2K"
//	carriageReturn := "\r"
//
//	// Build the full output string
//	var builder strings.Builder
//
//	builder.WriteString(saveCursor)      // Save cursor (input line)
//	builder.WriteString(moveCursorUp(1)) // Move up to output line
//
//	for _, line := range lines {
//		builder.WriteString(eraseLine)      // Clear line
//		builder.WriteString(carriageReturn) // Reset to beginning
//		builder.WriteString(line)           // Write output line
//		builder.WriteString("\n")           // Move to next line
//	}
//
//	builder.WriteString(restoreCursor) // Return to input line
//
//	// Send the whole update in one write
//	return []byte(builder.String())
//}

func CreateGameOutput(w *ecs.World, connectionId uuid.UUID, contents string) ecs.Entity {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(w, gameOutput, data.IsOutputComponent{})
	ecs.SetComponent(w, gameOutput, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(w, gameOutput, data.ContentsComponent{Contents: []byte(contents)})

	return gameOutput
}

func CreateClosingGameOutput(w *ecs.World, connectionId uuid.UUID, contents []byte) ecs.Entity {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(w, gameOutput, data.IsOutputComponent{})
	ecs.SetComponent(w, gameOutput, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(w, gameOutput, data.ContentsComponent{Contents: contents})
	ecs.SetComponent(w, gameOutput, data.CloseConnectionComponent{})

	return gameOutput
}
