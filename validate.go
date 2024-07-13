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
	entitySystemRegexp = regexp.MustCompile(`(?i)<!ENTITY\b`)
	doctypeBytes       = []byte(`DOCTYPE`)
)

// ValidateReader validates svg data from an io.Reader interface
func (vld Validator) ValidateReader(r io.Reader) error {
	t := xml.NewDecoder(r)
	var (
		to    xml.Token
		err   error
		elem  string
		id    string
		id4El string
		usec  = map[string]*useRef{}
		root  = &useRef{}
	)

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
			parent, ok := usec[id]
			if !ok {
				parent = root
			}
			var _id string
			_id, err = validateAttributes(v.Attr, vld.whiteListAttributes, vld.attrValueValidator, usec, parent)
			if err != nil {
				return err
			}
			if len(_id) > 0 {
				id = _id
				id4El = elem
			}
		case xml.EndElement:
			if id4El == elem {
				id = ``
			}
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

type useRef struct {
	parent *useRef
	count  uint64
	id     string
}

const maxReferences uint64 = 50

func (u *useRef) Add(n uint64) error {
	u.count += n
	if u.count > maxReferences {
		return fmt.Errorf(`%w: more than %d`, ErrTooManyReferences, maxReferences)
	}
	if u.parent != nil {
		return u.parent.Add(n)
	}
	return nil
}

func (u *useRef) String() string {
	return fmt.Sprintf(`{"id":%q,"count":%d,"parent":%s}`, u.id, u.count, u.parent)
}

func validateAttributes(attrs []xml.Attr, whiteListAttributes map[string]struct{}, attrValueValidator map[string]func(string) error, usec map[string]*useRef, parent *useRef) (id string, err error) {
	var key string
	for _, attr := range attrs {
		if len(attr.Name.Space) > 0 {
			switch attr.Name.Space {
			case "http://www.w3.org/XML/1998/namespace":
				attr.Name.Space = "xml"
			case "http://www.w3.org/1999/xlink":
				attr.Name.Space = "xlink"
			}
			key = strings.ToLower(attr.Name.Local)
			fn, ok := attrValueValidator[key]
			if ok {
				if err = fn(attr.Value); err != nil {
					return
				}
			}
			key = strings.ToLower(attr.Name.Space) + ":" + key
			if strings.HasSuffix(key, `xlink:href`) && strings.HasPrefix(attr.Value, `#`) {
				uk := strings.TrimPrefix(attr.Value, `#`)
				ref, ok := usec[uk]
				if !ok {
					ref = &useRef{parent: parent, id: uk}
					usec[uk] = ref
				}
				if err = ref.Add(1); err != nil {
					return
				}
			}
		} else {
			key = strings.ToLower(attr.Name.Local)
			if key == `id` {
				id = attr.Value
				_, ok := usec[id]
				if !ok {
					usec[id] = &useRef{parent: parent, id: id}
				}
			}
		}
		_, found := whiteListAttributes[key]
		if !found {
			err = fmt.Errorf("%w: %s", ErrInvalidAttribute, key)
			return
		}
		fn, ok := attrValueValidator[key]
		if ok {
			err = fn(attr.Value)
		} else {
			err = validateAttrValue(attr.Value)
		}
		if err != nil {
			return
		}
	}
	return
}

func validateElements(elm string, whiteListElements map[string]struct{}) bool {
	_, found := whiteListElements[elm]
	return found
}
