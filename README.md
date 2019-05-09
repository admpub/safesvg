# Safesvg
A Go library that will check if a given svg file is safe based on a whitelist of elements and attributes. This library does not sanitize svg files.

#### Word of caution
Using unsafe svg can be extremely dangerous. This library will not mitigate that risk. Please do your own research about svg security and risks before using this library.  

## Usage
```go
svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" d="M0 0h24v24H0V0z"/><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"/></svg>`)

v := safesvg.NewValidator()
err := v.ValidateBytes(validSvg)
if err != nil {
	fmt.Printf("Validation error %v", err)
}

```

Whitelist elements and attributes (adding to existing list, see validate.go)
```go
svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><newelement foo="bar" stranger="things"></newelement><path fill="none" d="M0 0h24v24H0V0z"/><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"/></svg>`)

v := safesvg.NewValidator()
v.WhitelistElements("newelement")
v.WhitelistAttributes("stranger","foo")

err := v.Validate(validSvg)
if err != nil {
	fmt.Printf("Validation error %v", err)
}
```

Blacklist elements and attributes (removing from existing list, see validate.go)
```go
svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" d="M0 0h24v24H0V0z"/><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"/></svg>`)

v := safesvg.NewValidator()
v.BlacklistElements("path")
v.BlacklistAttributes("width")

err := v.Validate(validSvg)
if err != nil {
	fmt.Printf("Validation error %v", err)
}
```

### Credits
The whitelist is copied from https://github.com/cure53/DOMPurify