package qbittorrent

import (
	"context"
	"net/url"
	"strings"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

// GetCategories returns every torrent category, keyed by category name.
func (c *Client) GetCategories(ctx context.Context) (map[string]*models.Category, error) {
	cats := map[string]*models.Category{}
	if err := c.getRequest(ctx, buildPath("/torrents/categories"), nil, &cats); err != nil {
		return nil, err
	}

	return cats, nil
}

// CreateCategory creates a new category.
func (c *Client) CreateCategory(ctx context.Context, req models.CreateCategoryRequest) error {
	params := url.Values{
		"category": {req.Name},
		"savePath": {req.SavePath},
	}

	return c.postRequest(ctx, buildPath("/torrents/createCategory"), params, nil)
}

// EditCategory changes the save path of an existing category.
func (c *Client) EditCategory(ctx context.Context, req models.EditCategoryRequest) error {
	params := url.Values{
		"category": {req.Name},
		"savePath": {req.SavePath},
	}

	return c.postRequest(ctx, buildPath("/torrents/editCategory"), params, nil)
}

// RemoveCategories removes the given categories. A request with no names is a
// no-op.
func (c *Client) RemoveCategories(ctx context.Context, req models.RemoveCategoriesRequest) error {
	if len(req.Names) == 0 {
		return nil
	}

	params := url.Values{
		"categories": {strings.Join(req.Names, "\n")},
	}

	return c.postRequest(ctx, buildPath("/torrents/removeCategories"), params, nil)
}
