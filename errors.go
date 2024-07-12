package safesvg

import "errors"

var (
	ErrInvalidElement              = errors.New("invalid element")
	ErrInvalidAttribute            = errors.New("invalid attribute")
	ErrUnallowedCSSAttributeValue  = errors.New("unallowed css attribute value")
	ErrUnallowedCSSAttribute       = errors.New("unallowed css attribute")
	ErrUnallowedHrefAttributeValue = errors.New("unallowed href attribute value")
)
