package lexer

type LexerMode uint8

const (
	ModeInitial LexerMode = iota
	ModeScripting
	ModeLookingForProperty
	ModeDoubleQuotes
	ModeNowDoc
	ModeHereDoc
	ModeEndHereDoc
	ModeBacktick
	ModeVarOffset
	ModeLookingForVarName
	ModeDocumentBlock
)

var /* const */ modeStrings = []string{
	"ModeInitial",
	"ModeScripting",
	"ModeLookingForProperty",
	"ModeDoubleQuotes",
	"ModeNowDoc",
	"ModeHereDoc",
	"ModeEndHereDoc",
	"ModeBacktick",
	"ModeVarOffset",
	"ModeLookingForVarName",
	"ModeDocumentBlock",
}

func (mode LexerMode) String() string {
	if int(mode) >= len(modeStrings) {
		return "Unknown"
	}
	return modeStrings[int(mode)]
}
