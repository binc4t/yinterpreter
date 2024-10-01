package identify

const (
	IDENT = "IDENT"
	INT   = "INT"
	INVALID

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	FOR      = "FOR"
	RETURN   = "RETURN"
	EOF      = "EOF"

	OPAssign   = "="
	OPPlus     = "+"
	OPMinus    = "-"
	OPAsterisk = "*"
	OPSlash    = "/"
	OPGT       = ">"
	OPLT       = "<"
	OPBang     = "!"
	OPEQ       = "=="
	OPNEQ      = "!="

	LParen    = "("
	RParen    = ")"
	LBrace    = "{"
	RBrace    = "}"
	Comma     = ","
	Semicolon = ";"
)

var keywords = map[string]string{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"for":    FOR,
	"return": RETURN,

	"=":  OPEQ,
	"+":  OPPlus,
	"-":  OPMinus,
	"*":  OPAsterisk,
	"/":  OPSlash,
	">":  OPGT,
	"<":  OPLT,
	"==": OPEQ,
	"!=": OPNEQ,

	"(": LParen,
	")": RParen,
	"{": LBrace,
	"}": RBrace,
	",": Comma,
	";": Semicolon,
}

func LookupKeywords(s string) string {
	if k, ok := keywords[s]; ok {
		return k
	}
	return ""
}
