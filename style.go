package safesvg

import (
	"errors"
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
				return errors.New("Unallowed css attribute value " + token.Value)
			}
		case scanner.TokenAtKeyword:
			if strings.EqualFold(token.Value, `@import`) {
				return errors.New("Unallowed css attribute " + token.Value)
			}
		case scanner.TokenFunction: // expression(...) regex(...)
			return errors.New("Unallowed css attribute value " + token.Value)
		}
	}
	return err
}
