package models

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

// EditCategoryRequest changes the save path of an existing category.
type EditCategoryRequest struct {
	// Name is the existing category name. Required.
	Name string
	// SavePath is the new save path.
	SavePath string
}

// RemoveCategoriesRequest removes one or more categories by name.
type RemoveCategoriesRequest struct {
	Names []string
}
