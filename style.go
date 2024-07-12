package safesvg

import (
	"fmt"
	"strings"

	"github.com/gorilla/css/scanner"
)

var cssDebug bool

func ValidateStyle(myCSS []byte) error {
	var err error
	s := scanner.New(string(myCSS))
	for {
		token := s.Next()
		if token.Type == scanner.TokenEOF || token.Type == scanner.TokenError {
			break
		}
		// Do something with the token...
		if cssDebug {
			println(token.Type.String(), `=====>`, token.Value)
		}
		switch token.Type {
		case scanner.TokenURI:
			if strings.Contains(token.Value, `//`) { // url(...)
				return fmt.Errorf("%w: %s", ErrUnallowedCSSAttributeValue, token.Value)
			}
		case scanner.TokenAtKeyword:
			if strings.EqualFold(token.Value, `@import`) {
				return fmt.Errorf("%w: %s", ErrUnallowedCSSAttribute, token.Value)
			}
		case scanner.TokenFunction: // expression(...) regex(...)
			return fmt.Errorf("%w: %s", ErrUnallowedCSSAttributeValue, token.Value)
		}
	}
	return err
}
