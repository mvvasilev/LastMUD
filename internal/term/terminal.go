package term

import (
	"bufio"
	"bytes"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"context"
	"encoding/binary"
	"fmt"
	"slices"
	"sync"
	"time"
)

type termError struct {
	err string
}

func createTermError(v ...any) *termError {
	return &termError{
		err: fmt.Sprint(v...),
	}
}

func (te *termError) Error() string {
	return te.err
}

// NAWS ( RFC 1073 ):
//   - Server -> Client: IAC DO NAWS
//   - Client -> Server: IAC WILL NAWS
//   - Client -> Server: IAC SB NAWS 0 <width> 0 <height> IAC SE
//
// ECHO ( RFC 857 )
//   - Server -> Client: IAC WILL ECHO
//   - Client -> Server: IAC DO ECHO
//
// SUPPRESSGOAHEAD ( RFC 858 )
//   - Server -> Client: IAC WILL SUPPRESSGOAHEAD
//   - Client -> Server: IAC DO SUPPRESSGOAHEAD
//
// LINEMODE ( RFC 1184 )
//   - Server -> Client: IAC DONT LINEMODE
//   - Client -> Server: IAC WONT LINEMODE
const (
	IAC             byte = 255
	DO              byte = 253
	DONT            byte = 254
	WILL            byte = 251
	WONT            byte = 252
	SB              byte = 250
	SE              byte = 240
	ECHO            byte = 1
	SUPPRESSGOAHEAD byte = 3
	NAWS            byte = 31
	LINEMODE        byte = 34
)

func telnetByteToString(telnetByte byte) string {
	switch telnetByte {
	case IAC:
		return "IAC"
	case DO:
		return "DO"
	case DONT:
		return "DONT"
	case WILL:
		return "WILL"
	case WONT:
		return "WONT"
	case SB:
		return "SB"
	case SE:
		return "SE"
	case ECHO:
		return "ECHO"
	case SUPPRESSGOAHEAD:
		return "SUPPRESS-GO-AHEAD"
	case NAWS:
		return "NAWS"
	case LINEMODE:
		return "LINEMODE"
	default:
		return fmt.Sprintf("%02X", telnetByte)
	}
}

const (
	ClearScreen           = "\x1b[2J"
	ClearLine             = "\x1b[K"
	MoveCursorStartOfLine = "\x1b[G"
	SaveCursorPosition    = "\x1b[s"
	RestoreCursorPosition = "\x1b[u"
	ScrollUp              = "\x1b[1S"
	MoveCursorUpOne       = "\x1b[1A"
	MoveCursorRightOne    = "\x1b[1C"
	MoveCursorLeftOne     = "\x1b[1D"
)

type VirtualTerm struct {
	ctx context.Context
	wg  *sync.WaitGroup

	buffer *bytes.Buffer
	cX, xY int

	width, height int

	timeout func(t time.Time)

	reader *bufio.Reader
	writer *bufio.Writer

	writes    chan []byte
	submitted chan string

	stop context.CancelFunc
}

func CreateVirtualTerm(ctx context.Context, wg *sync.WaitGroup, timeout func(t time.Time), reader *bufio.Reader, writer *bufio.Writer) (term *VirtualTerm, err error) {
	ctx, cancel := context.WithCancel(ctx)
	term = &VirtualTerm{
		ctx:       ctx,
		stop:      cancel,
		wg:        wg,
		buffer:    bytes.NewBuffer([]byte{}),
		cX:        0,
		xY:        0,
		width:     80, // Default of 80
		height:    24, // Default of 24
		timeout:   timeout,
		reader:    reader,
		writer:    writer,
		submitted: make(chan string, 1),
		writes:    make(chan []byte, 100),
	}

	err = term.sendWillSuppressGA()
	err = term.sendWillEcho()
	err = term.sendDisableLinemode()
	err = term.sendNAWSNegotiationRequest()

	if err != nil {
		logging.Error(err)
	}

	wg.Add(2)
	go term.listen()
	go term.send()

	return
}

func (term *VirtualTerm) Close() {
	term.stop()
}

func (term *VirtualTerm) NextCommand() (cmd string) {
	if term.shouldStop() {
		return
	}

	select {
	case cmd = <-term.submitted:
		return cmd
	default:
		return
	}
}

func (term *VirtualTerm) Write(bytes []byte) (err error) {
	if term.shouldStop() {
		return
	}

	term.writes <- bytes
	return
}

func (term *VirtualTerm) readByte() (b byte, err error) {
	term.timeout(time.Now().Add(1000 * time.Millisecond))
	b, err = term.reader.ReadByte()
	return
}

func (term *VirtualTerm) listen() {
	defer term.wg.Done()

	for {
		if term.shouldStop() {
			break
		}

		b, err := term.readByte()

		if err != nil {
			continue
		}

		switch b {
		case IAC:
			err = term.handleTelnetCommand()

			if err != nil {
				term.Close()
				break
			}

			continue
		case '\b', 127:
			if term.buffer.Len() <= 0 {
				continue
			}

			term.buffer = bytes.NewBuffer(term.buffer.Bytes()[0 : len(term.buffer.Bytes())-1])
			term.writer.Write([]byte("\x1b[D"))
			term.writer.Write([]byte("\x1b[P"))
		case '\r', '\n':
			if !term.shouldStop() {
				term.submitted <- term.buffer.String()
			}
			term.buffer = bytes.NewBuffer([]byte{})
		case '\t', '\a', '\f', '\v':
			continue
		default:
		}

		if isInputCharacter(b) {
			term.buffer.WriteByte(b)
		}

		term.writer.Write([]byte(ClearLine))
		term.writer.Write([]byte(MoveCursorStartOfLine))
		term.writer.Write([]byte("> "))
		term.writer.Write(term.buffer.Bytes())
	}
}

func isInputCharacter(b byte) bool {
	return b >= 0x20 && b <= 0x7E
}

func (term *VirtualTerm) send() {
	defer term.wg.Done()

	_ = term.sendClear()

	for {
		if term.shouldStop() {
			break
		}

		select {
		case write := <-term.writes:
			for _, w := range term.formatOutput(write) {
				term.writer.Write([]byte(MoveCursorStartOfLine))
				term.writer.Write([]byte(ClearLine))
				term.writer.Write([]byte(w))
				term.writer.Write([]byte("\r\n"))
			}

			term.writer.Write([]byte("> "))
			term.writer.Write(term.buffer.Bytes())
		default:
		}

		err := term.writer.Flush()

		if err != nil {
			logging.Error(err)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (term *VirtualTerm) sendNAWSNegotiationRequest() (err error) {
	_, err = term.writer.Write([]byte{IAC, DO, NAWS})
	return
}

func (term *VirtualTerm) sendWillEcho() (err error) {
	_, err = term.writer.Write([]byte{IAC, WILL, ECHO})
	return
}

func (term *VirtualTerm) sendWillSuppressGA() (err error) {
	_, err = term.writer.Write([]byte{IAC, WILL, SUPPRESSGOAHEAD})
	return
}

func (term *VirtualTerm) sendDisableLinemode() (err error) {
	_, err = term.writer.Write([]byte{IAC, DONT, LINEMODE})
	return
}

func (term *VirtualTerm) sendClear() (err error) {
	_, err = term.writer.Write([]byte(ClearScreen))
	return
}

func (term *VirtualTerm) handleTelnetCommand() (err error) {
	buf := make([]byte, 255)
	buf[0] = IAC

	next, err := term.readByte()

	if err != nil {
		return err
	}

	lastIndex := 1

	switch next {
	case IAC:
		// Double IAC, meant to be interpreted as literal
		term.buffer.WriteByte(255)
	case DO, DONT, WILL, WONT:
		// Negotiation
		lastIndex++
		buf[lastIndex] = next
		final, err := term.readByte()

		if err != nil {
			return err
		}

		lastIndex++
		buf[lastIndex] = final
	default:
		// Send begin, send end
		for {
			next, err := term.readByte()

			if err != nil {
				return err
			}

			buf[lastIndex] = next

			if next == SE {
				break
			}

			lastIndex++
		}
	}

	strRep := make([]string, lastIndex+1)

	for i := 0; i < lastIndex+1; i++ {
		strRep[i] = telnetByteToString(buf[i])
	}

	if slices.Equal([]byte{IAC, WONT, NAWS}, buf[:3]) {
		// Client does not agree to NAWS, cannot proceed
		return createTermError("NAWS negotiation failed")
	} else if slices.Equal([]byte{IAC, DONT, SUPPRESSGOAHEAD}, buf[:3]) {
		// Client does not agree to suppress go-ahead
		return createTermError("suppress-go-ahead negotiation failed")
	} else if slices.Equal([]byte{IAC, DONT, ECHO}, buf[:3]) {
		// Client does not agree to not echo
		return createTermError("No echo negotiation failed")
	} else if slices.Equal([]byte{IAC, WILL, LINEMODE}, buf[:3]) {
		return createTermError("Client wants to use linemode")
	} else if slices.Equal([]byte{IAC, SB, NAWS}, buf[:3]) {
		logging.Info("Received NAWS Response")
		// Client sending NAWS data
		term.width = int(binary.BigEndian.Uint16(buf[3:5]))
		term.height = int(binary.BigEndian.Uint16(buf[5:7]))
	}

	return
}

func (term *VirtualTerm) shouldStop() bool {
	select {
	case <-term.ctx.Done():
		return true
	default:
	}

	return false
}

func (term *VirtualTerm) formatOutput(output []byte) []string {
	return []string{string(output)}

	// TODO
	//strText := string(output)
	//
	//strText = strings.ReplaceAll(strText, "\n", " ")
	//words := strings.Fields(strText)
	//
	//var lines [][]string
	//
	//for _, word := range words {
	//
	//	//if len(line)+len(word) > term.width {
	//	//	lines = append(lines, line)
	//	//} else {
	//	//	line += " " + word
	//	//}
	//}
	//
	//return lines
}
