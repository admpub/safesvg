package safesvg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Validator is a struct with private variables for storing the whitelists
type Validator struct {
	whiteListElements   map[string]struct{}
	whiteListAttributes map[string]struct{}
	innerTextValidator  map[string]func([]byte) error
	attrValueValidator  map[string]func(string) error
}

// NewValidator creates a new validator with default whitelists
func NewValidator() Validator {
	vld := Validator{
		whiteListElements:   map[string]struct{}{},
		whiteListAttributes: map[string]struct{}{},
		innerTextValidator: map[string]func([]byte) error{
			`style`: ValidateStyle,
		},
		attrValueValidator: map[string]func(string) error{
			`href`: validateHref,
		},
	}
	for k, v := range svg_elements {
		vld.whiteListElements[k] = v
	}
	for k, v := range svg_attributes {
		vld.whiteListAttributes[k] = v
	}
	return vld
}

// Validate validates a slice of bytes containing the svg data
func (vld Validator) Validate(b []byte) error {
	r := bytes.NewReader(b)
	return vld.ValidateReader(r)
}

var (
	entitySystemRegexp = regexp.MustCompile(`(?i)<!ENTITY(?:.*)SYSTEM`)
	doctypeBytes       = []byte(`DOCTYPE`)
)

// ValidateReader validates svg data from an io.Reader interface
func (vld Validator) ValidateReader(r io.Reader) error {
	t := xml.NewDecoder(r)
	var to xml.Token
	var err error
	var elem string

	for {
		to, err = t.Token()
		if err != nil {
			if err == io.EOF || err.Error() == "EOF" {
				break
			}
			return err
		}

		switch v := to.(type) {
		case xml.StartElement:
			elem = strings.ToLower(v.Name.Local)
			if ok := validateElements(elem, vld.whiteListElements); !ok {
				return fmt.Errorf("%w: %s", ErrInvalidElement, v.Name.Local)
			}

			if err = validateAttributes(v.Attr, vld.whiteListAttributes, vld.attrValueValidator); err != nil {
				return err
			}
		case xml.EndElement:
			elem = ``
			if ok := validateElements(strings.ToLower(v.Name.Local), vld.whiteListElements); !ok {
				return fmt.Errorf("%w: %s", ErrInvalidElement, v.Name.Local)
			}
		case xml.CharData: //text
			if len(elem) > 0 {
				if fn, ok := vld.innerTextValidator[elem]; ok {
					if err := fn(v); err != nil {
						return err
					}
				}
			}

		case xml.Comment: // <!--...-->

		case xml.ProcInst: // <?target inst?>
			if !strings.EqualFold(v.Target, `xml`) {
				return fmt.Errorf("%w: %s", ErrInvalidElement, v.Target)
			}

		case xml.Directive: // <!...> doctype etc
			d := bytes.TrimSpace(v)
			length := len(d)
			if length > 8 && bytes.EqualFold(d[0:7], doctypeBytes) {
				if entitySystemRegexp.Match(d) {
					return fmt.Errorf("%w: %s", ErrUnallowedEntityAttribute, d)
				}
			}
		}

	}

	return nil
}

// WhitelistElements adds svg elements to the whitelist
func (vld *Validator) WhitelistElements(elements ...string) *Validator {
	for _, elemet := range elements {
		elemet = strings.ToLower(elemet)
		vld.whiteListElements[elemet] = struct{}{}
	}
	return vld
}

// WhitelistAttributes adds svg attributes to the whitelist
func (vld *Validator) WhitelistAttributes(attributes ...string) *Validator {
	for _, attr := range attributes {
		attr = strings.ToLower(attr)
		vld.whiteListAttributes[attr] = struct{}{}
	}
	return vld
}

// BlacklistElements removes svg elements from the whitelist
func (vld *Validator) BlacklistElements(elements ...string) *Validator {
	for _, elemet := range elements {
		elemet = strings.ToLower(elemet)
		delete(vld.whiteListElements, elemet)
	}
	return vld
}

// BlacklistAttributes removes svg attributes from the whitelist
func (vld *Validator) BlacklistAttributes(attributes ...string) *Validator {
	for _, attr := range attributes {
		attr = strings.ToLower(attr)
		delete(vld.whiteListAttributes, attr)
	}
	return vld
}

func (vld *Validator) SetInnerTextValidator(element string, validate func([]byte) error) *Validator {
	element = strings.ToLower(element)
	vld.innerTextValidator[element] = validate
	return vld
}

func (vld *Validator) SetAttrValueValidator(attribute string, validate func(string) error) *Validator {
	attribute = strings.ToLower(attribute)
	vld.attrValueValidator[attribute] = validate
	return vld
}

func (vld *Validator) RemoveInnerTextValidator(element string) *Validator {
	element = strings.ToLower(element)
	delete(vld.innerTextValidator, element)
	return vld
}

func (vld *Validator) RemoveAttrValueValidator(attribute string) *Validator {
	attribute = strings.ToLower(attribute)
	delete(vld.attrValueValidator, attribute)
	return vld
}

func validateAttributes(attrs []xml.Attr, whiteListAttributes map[string]struct{}, attrValueValidator map[string]func(string) error) error {
	var key string
	var err error
	for _, attr := range attrs {
		if len(attr.Name.Space) > 0 {
			if attr.Name.Space == "http://www.w3.org/XML/1998/namespace" {
				attr.Name.Space = "xml"
			}
			key = strings.ToLower(attr.Name.Local)
			fn, ok := attrValueValidator[key]
			if ok {
				if err = fn(attr.Value); err != nil {
					return err
				}
			}
			key = strings.ToLower(attr.Name.Space) + ":" + key
		} else {
			key = strings.ToLower(attr.Name.Local)
		}
		_, found := whiteListAttributes[key]
		if !found {
			return fmt.Errorf("%w: %s", ErrInvalidAttribute, key)
		}
		fn, ok := attrValueValidator[key]
		if ok {
			err = fn(attr.Value)
		} else {
			err = validateAttrValue(attr.Value)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func validateElements(elm string, whiteListElements map[string]struct{}) bool {
	_, found := whiteListElements[elm]
	return found
}
