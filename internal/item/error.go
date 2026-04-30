package item

import "errors"

var ErrItemNotFound = errors.New("Item not found")
var ErrAccessDenied = errors.New("Access denied")
var ErrAlreadyReserved = errors.New("item already reserved")
