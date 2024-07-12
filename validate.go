package safesvg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

var svg_elements = map[string]struct{}{
	"svg":                 {},
	"altglyph":            {},
	"altglyphdef":         {},
	"altglyphitem":        {},
	"animatecolor":        {},
	"animatemotion":       {},
	"animatetransform":    {},
	"circle":              {},
	"clippath":            {},
	"defs":                {},
	"desc":                {},
	"ellipse":             {},
	"filter":              {},
	"font":                {},
	"g":                   {},
	"glyph":               {},
	"glyphref":            {},
	"hkern":               {},
	"image":               {},
	"line":                {},
	"lineargradient":      {},
	"marker":              {},
	"mask":                {},
	"metadata":            {},
	"mpath":               {},
	"path":                {},
	"pattern":             {},
	"polygon":             {},
	"polyline":            {},
	"radialgradient":      {},
	"rect":                {},
	"stop":                {},
	"switch":              {},
	"symbol":              {},
	"text":                {},
	"textpath":            {},
	"title":               {},
	"tref":                {},
	"tspan":               {},
	"use":                 {},
	"view":                {},
	"vkern":               {},
	"feblend":             {},
	"fecolormatrix":       {},
	"fecomponenttransfer": {},
	"fecomposite":         {},
	"feconvolvematrix":    {},
	"fediffuselighting":   {},
	"fedisplacementmap":   {},
	"fedistantlight":      {},
	"feflood":             {},
	"fefunca":             {},
	"fefuncb":             {},
	"fefuncg":             {},
	"fefuncr":             {},
	"fegaussianblur":      {},
	"femerge":             {},
	"femergenode":         {},
	"femorphology":        {},
	"feoffset":            {},
	"fepointlight":        {},
	"fespecularlighting":  {},
	"fespotlight":         {},
	"fetile":              {},
	"feturbulence":        {},
}

var svg_attributes = map[string]struct{}{
	"accent-height":               {},
	"accumulate":                  {},
	"additivive":                  {},
	"alignment-baseline":          {},
	"ascent":                      {},
	"attributename":               {},
	"attributetype":               {},
	"azimuth":                     {},
	"baseprofile":                 {},
	"basefrequency":               {},
	"baseline-shift":              {},
	"begin":                       {},
	"bias":                        {},
	"by":                          {},
	"class":                       {},
	"clip":                        {},
	"clip-path":                   {},
	"clip-rule":                   {},
	"color":                       {},
	"color-interpolation":         {},
	"color-interpolation-filters": {},
	"color-profile":               {},
	"color-rendering":             {},
	"cx":                          {},
	"cy":                          {},
	"d":                           {},
	"dx":                          {},
	"dy":                          {},
	"diffuseconstant":             {},
	"direction":                   {},
	"display":                     {},
	"divisor":                     {},
	"dur":                         {},
	"edgemode":                    {},
	"elevation":                   {},
	"end":                         {},
	"fill":                        {},
	"fill-opacity":                {},
	"fill-rule":                   {},
	"filter":                      {},
	"flood-color":                 {},
	"flood-opacity":               {},
	"font-family":                 {},
	"font-size":                   {},
	"font-size-adjust":            {},
	"font-stretch":                {},
	"font-style":                  {},
	"font-variant":                {},
	"font-weight":                 {},
	"fx":                          {},
	"fy":                          {},
	"g1":                          {},
	"g2":                          {},
	"glyph-name":                  {},
	"glyphref":                    {},
	"gradientunits":               {},
	"gradienttransform":           {},
	"height":                      {},
	"href":                        {},
	"id":                          {},
	"image-rendering":             {},
	"in":                          {},
	"in2":                         {},
	"k":                           {},
	"k1":                          {},
	"k2":                          {},
	"k3":                          {},
	"k4":                          {},
	"kerning":                     {},
	"keypoints":                   {},
	"keysplines":                  {},
	"keytimes":                    {},
	"lang":                        {},
	"lengthadjust":                {},
	"letter-spacing":              {},
	"kernelmatrix":                {},
	"kernelunitlength":            {},
	"lighting-color":              {},
	"local":                       {},
	"marker-end":                  {},
	"marker-mid":                  {},
	"marker-start":                {},
	"markerheight":                {},
	"markerunits":                 {},
	"markerwidth":                 {},
	"maskcontentunits":            {},
	"maskunits":                   {},
	"max":                         {},
	"mask":                        {},
	"media":                       {},
	"method":                      {},
	"mode":                        {},
	"min":                         {},
	"name":                        {},
	"numoctaves":                  {},
	"offset":                      {},
	"operator":                    {},
	"opacity":                     {},
	"order":                       {},
	"orient":                      {},
	"orientation":                 {},
	"origin":                      {},
	"overflow":                    {},
	"paint-order":                 {},
	"path":                        {},
	"pathlength":                  {},
	"patterncontentunits":         {},
	"patterntransform":            {},
	"patternunits":                {},
	"points":                      {},
	"preservealpha":               {},
	"preserveaspectratio":         {},
	"r":                           {},
	"rx":                          {},
	"ry":                          {},
	"radius":                      {},
	"refx":                        {},
	"refy":                        {},
	"repeatcount":                 {},
	"repeatdur":                   {},
	"restart":                     {},
	"result":                      {},
	"rotate":                      {},
	"scale":                       {},
	"seed":                        {},
	"shape-rendering":             {},
	"specularconstant":            {},
	"specularexponent":            {},
	"spreadmethod":                {},
	"stddeviation":                {},
	"stitchtiles":                 {},
	"stop-color":                  {},
	"stop-opacity":                {},
	"stroke-dasharray":            {},
	"stroke-dashoffset":           {},
	"stroke-linecap":              {},
	"stroke-linejoin":             {},
	"stroke-miterlimit":           {},
	"stroke-opacity":              {},
	"stroke":                      {},
	"stroke-width":                {},
	"style":                       {},
	"surfacescale":                {},
	"tabindex":                    {},
	"targetx":                     {},
	"targety":                     {},
	"transform":                   {},
	"text-anchor":                 {},
	"text-decoration":             {},
	"text-rendering":              {},
	"textlength":                  {},
	"type":                        {},
	"u1":                          {},
	"u2":                          {},
	"unicode":                     {},
	"version":                     {},
	"values":                      {},
	"viewbox":                     {},
	"visibility":                  {},
	"vert-adv-y":                  {},
	"vert-origin-x":               {},
	"vert-origin-y":               {},
	"width":                       {},
	"word-spacing":                {},
	"wrap":                        {},
	"writing-mode":                {},
	"xchannelselector":            {},
	"ychannelselector":            {},
	"x":                           {},
	"x1":                          {},
	"x2":                          {},
	"xmlns":                       {},
	"y":                           {},
	"y1":                          {},
	"y2":                          {},
	"z":                           {},
	"zoomandpan":                  {},

	"xlink:href":  {},
	"xml:id":      {},
	"xlink:title": {},
	"xml:space":   {},
	"xmlns:xlink": {},
}

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
		whiteListElements:   svg_elements,
		whiteListAttributes: svg_attributes,
		innerTextValidator: map[string]func([]byte) error{
			`style`: ValidateStyle,
		},
		attrValueValidator: map[string]func(string) error{},
	}
	return vld
}

// Validate validates a slice of bytes containing the svg data
func (vld Validator) Validate(b []byte) error {
	r := bytes.NewReader(b)
	return vld.ValidateReader(r)
}

// ValidateReader validates svg data from an io.Reader interface
func (vld Validator) ValidateReader(r io.Reader) error {
	t := xml.NewDecoder(r)
	var to xml.Token
	var err error
	var elem string

	for {
		to, err = t.Token()

		switch v := to.(type) {
		case xml.StartElement:
			elem = strings.ToLower(v.Name.Local)
			if ok := validElements(elem, vld.whiteListElements); !ok {
				return fmt.Errorf("%w: %s", ErrInvalidElement, v.Name.Local)
			}

			if err := validAttributes(v.Attr, vld.whiteListAttributes, vld.attrValueValidator); err != nil {
				return err
			}
		case xml.EndElement:
			elem = ``
			if ok := validElements(strings.ToLower(v.Name.Local), vld.whiteListElements); !ok {
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

		case xml.Comment:

		case xml.ProcInst:

		case xml.Directive: //doctype etc

		}

		if err != nil {
			if err == io.EOF || err.Error() == "EOF" {
				break
			}
			return err
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

func (vld *Validator) RemoveInnerTextValidator(element string, validate func([]byte) error) *Validator {
	element = strings.ToLower(element)
	delete(vld.innerTextValidator, element)
	return vld
}

func (vld *Validator) RemoveAttrValueValidator(attribute string, validate func(string) error) *Validator {
	attribute = strings.ToLower(attribute)
	delete(vld.attrValueValidator, attribute)
	return vld
}

func validAttributes(attrs []xml.Attr, whiteListAttributes map[string]struct{}, attrValueValidator map[string]func(string) error) error {
	var key string
	var err error
	for _, attr := range attrs {
		if len(attr.Name.Space) > 0 {
			if attr.Name.Space == "http://www.w3.org/XML/1998/namespace" {
				attr.Name.Space = "xml"
			}
			key = attr.Name.Space + ":" + attr.Name.Local
		} else {
			key = attr.Name.Local
		}
		key = strings.ToLower(key)
		_, found := whiteListAttributes[key]
		if !found {
			return fmt.Errorf("%w: %s", ErrInvalidAttribute, attr.Name.Local)
		}
		fn, ok := attrValueValidator[key]
		if ok {
			if err = fn(attr.Value); err != nil {
				return err
			}
		}
	}
	return err
}

func validElements(elm string, whiteListElements map[string]struct{}) bool {
	_, found := whiteListElements[elm]
	return found
}
