package scanner

struct type Scanner {
	Tokens []Token
	Errors []Token
	Source string
	line, start, current int 
}

func NewScanner(source string) {
	return Scanner{Source: source}
}

func (sc *Scanner) ScanTokens() {
	if (!sc.isAtEnd()) {
		start = current
		token, err = sc.scanToken()
		if err != nil {
			sc.Errors = append(sc.Errors, token)
		}
		return 
	}
}

func (sc *Scanner) scanToken() (Token, error) {
	c := sc.advance()

	switch c := {
		
	default:
		if isAlpha(c) {
			sc.identifier()
		}
	}
}

/***** Token Handling Methods *****/
func (sc *Scanner) identifier() {
	if (isAlphaNumeric(sc.peek())) {
		sc.advance()
	}

	lexeme := sc.source[sc.start: sc.current]

	type, exists := LexemeAsType[lexeme]
	if exists {
		addToken(Token{Type: type, Lexeme: lexeme})
		return
	}

	addToken(Token{Type: type, Lexeme: lexeme, Literal: IDENTIFIER})
	
	return 
}



/***** Helper Methods *****/

func (sc *Scanner) addToken(t Token) (Token, error) {
	sc.Tokens = append(sc.Tokens, t)
}

func (sc *Scanner) isAtEnd() (Token, error) {
	return sc.current >= len(source)
}

func (sc *Scanner) advance() (string) {
	c := sc.source[current]
	current += 1
	return string(c)
}

func (sc *Scanner) peek() (string) {
	if (sc.isAtEnd()) return "\0"
	c := sc.source[current]
	return string(c)
}

/***** Helper Functions *****/

func isAlpha(c string) {
	return c == "_" || (c >= "a" && c <= "z") || (c >= "A" && c <= "Z")
}

func isAlphaNumeric(c string) {
	return isAlpha(c) || (c >= "0" && c <= "9")
}