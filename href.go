package safesvg

import (
	"fmt"
	"regexp"
	"strings"
)

var hrefDataRegex = regexp.MustCompile(`(?i)^\s*[^/]+/[^/;]+\s*;\s*`)
var hrefDataMimes = []string{`image/png`, `image/jpg`, `image/jpeg`, `image/pjpeg`, `image/gif`}

func validateHref(value string) error {
	value = strings.TrimSpace(value)
	if err := validateAttrValue(value); err != nil {
		return err
	}
	if len(value) > 5 && strings.EqualFold(value[0:5], `data:`) { // data:image/png;base64,
		value2 := value[5:]
		if !hrefDataRegex.MatchString(value2) {
			return nil
		}
		mime := strings.SplitN(value2, `;`, 2)[0]
		mime = strings.ToLower(mime)
		for _, allowed := range hrefDataMimes {
			if allowed == mime {
				return nil
			}
		}
		return fmt.Errorf(`%w: %s`, ErrUnallowedHrefAttributeValue, value)
	}
	return nil
}

func validateAttrValue(value string) error {
	value = strings.TrimSpace(value)
	length := len(value)
	switch {
	case length > 11:
		if strings.EqualFold(value[0:11], `javascript:`) {
			return fmt.Errorf(`%w: %s`, ErrUnallowedHrefAttributeValue, value)
		}
	case length == 11:
		if strings.EqualFold(value, `javascript:`) {
			return fmt.Errorf(`%w: %s`, ErrUnallowedHrefAttributeValue, value)
		}
	}
	return nil
}
