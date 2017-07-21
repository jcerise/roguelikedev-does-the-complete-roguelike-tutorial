package ui

import (
	blt "bearlibterminal"
)

type MessageLog struct {
	messages []string
	MaxLength int
}

func (ml *MessageLog) InitMessages(maxLength int) {
	ml.messages = make([]string, maxLength)
}

func (ml *MessageLog) SendMessage(message string) {
	// Prepend the message onto the messageLog slice
	if len(ml.messages) >= ml.MaxLength {
		// Throw away any messages that exceed our total queue size
		ml.messages = ml.messages[:len(ml.messages)-1]
	}
	ml.messages = append([]string{message}, ml.messages...)
}

func (ml *MessageLog) PrintMessages(viewAreaY, windowSizeX, windowSizeY int) {
	// Print the latest five messages from the messageLog. These will be printed in reverse order (newest at the top),
	// to make it appear they are scrolling down the screen
	clearMessages(viewAreaY, windowSizeX, windowSizeY)

	toShow := 0

	if len(ml.messages) <= 5 {
		// Just loop through the messageLog, printing them in reverse order
		toShow = len(ml.messages)
	} else {
		// If we have more than 5 messages stored, just show the five most recent
		toShow = 5
	}

	blt.Color(blt.ColorFromName("white"))
	blt.Layer(1)
	for i := toShow; i > 0; i-- {
		blt.Print(1, (viewAreaY - 1) + i, ml.messages[i - 1])
	}
}

func clearMessages(viewAreaY, windowSizeX, windowSizeY int) {
	// Clear the message area, so our messages do not overlap
	blt.ClearArea(0, viewAreaY, windowSizeX, windowSizeY - viewAreaY)
}