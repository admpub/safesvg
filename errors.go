package safesvg

import "errors"

var (
	ErrInvalidElement              = errors.New("[svg] invalid element")
	ErrInvalidAttribute            = errors.New("[svg] invalid attribute")
	ErrUnallowedCSSAttributeValue  = errors.New("[svg] unallowed css attribute value")
	ErrUnallowedCSSAttribute       = errors.New("[svg] unallowed css attribute")
	ErrUnallowedHrefAttributeValue = errors.New("[svg] unallowed href attribute value")
	ErrUnallowedEntityAttribute    = errors.New("[svg] unallowed entity attribute")
	ErrTooManyReferences           = errors.New("[svg] too many references")
)
