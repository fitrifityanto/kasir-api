package models

import "errors"

var ErrCategoryHasProducts = errors.New("category has products")
var ErrCategoryNotFound = errors.New("category not found")
