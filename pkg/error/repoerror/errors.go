package repoerror

import "fmt"

var (
	ErrNotFound = fmt.Errorf("not found")
	ErrExisted  = fmt.Errorf("existed")
)
