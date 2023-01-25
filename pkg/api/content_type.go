package api

import "fmt"

type ContentType int

const (
	ContentTypeUnknown ContentType = iota
	ContentTypeJSON
)

func (c ContentType) String() string {
	switch c {
	case ContentTypeJSON:
		return "application/json"
	case ContentTypeUnknown:
		return "unknown"
	}

	return ""
}

func DeriveContentType(accept string) ContentType {
	switch accept {
	case "application/json":
		return ContentTypeJSON
	case "*/*":
		return ContentTypeJSON
	case "":
		return ContentTypeJSON
	}

	return ContentTypeUnknown
}

func ValidateContentType(contentType ContentType, accepting []ContentType) error {
	if !DoesAccept(accepting, contentType) {
		return fmt.Errorf("unsupported content-type: %s", contentType.String())
	}

	return nil
}
