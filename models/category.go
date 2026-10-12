package models

import "net/url"

// Category is a torrent category as returned by the Web API.
type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// CreateCategoryRequest creates a new category.
type CreateCategoryRequest struct {
	// Name is the category name. Required.
	Name string
	// SavePath is the category save path. When empty, the default save path
	// is used.
	SavePath string
}

// ToValues encodes the request into URL-form parameters.
func (r *CreateCategoryRequest) ToValues() url.Values {
	values := url.Values{}
	values.Set("category", r.Name)
	values.Set("savePath", r.SavePath)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *CreateCategoryRequest) FromValues(values url.Values) error {
	r.Name = values.Get("category")
	r.SavePath = values.Get("savePath")

	return nil
}

// EditCategoryRequest changes the save path of an existing category.
type EditCategoryRequest struct {
	// Name is the existing category name. Required.
	Name string
	// SavePath is the new save path.
	SavePath string
}

// ToValues encodes the request into URL-form parameters.
func (r *EditCategoryRequest) ToValues() url.Values {
	values := url.Values{}
	values.Set("category", r.Name)
	values.Set("savePath", r.SavePath)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *EditCategoryRequest) FromValues(values url.Values) error {
	r.Name = values.Get("category")
	r.SavePath = values.Get("savePath")

	return nil
}

// RemoveCategoriesRequest removes one or more categories by name.
type RemoveCategoriesRequest struct {
	// Names are the category names to remove.
	Names []string
}

// ToValues encodes the request into URL-form parameters.
func (r *RemoveCategoriesRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "categories", r.Names, "\n")

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *RemoveCategoriesRequest) FromValues(values url.Values) error {
	r.Names = listParam(values, "categories", "\n")

	return nil
}
