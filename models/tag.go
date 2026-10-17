package models

import "net/url"

// Tag is a tag name as returned by the Web API.
type Tag struct {
	Name string `json:"name"`
}

// CreateTagRequest creates one or more tags.
type CreateTagRequest struct {
	// Tags are the tag names to create.
	Tags []string
}

// ToValues encodes the request into URL-form parameters.
func (r *CreateTagRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "tags", r.Tags, ",")

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *CreateTagRequest) FromValues(values url.Values) error {
	r.Tags = listParam(values, "tags", ",")

	return nil
}

// DeleteTagRequest deletes one or more tags.
type DeleteTagRequest struct {
	// Tags are the tag names to delete.
	Tags []string
}

// ToValues encodes the request into URL-form parameters.
func (r *DeleteTagRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "tags", r.Tags, ",")

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *DeleteTagRequest) FromValues(values url.Values) error {
	r.Tags = listParam(values, "tags", ",")

	return nil
}
