package router

import (
	"fmt"
	"strings"
)

type URLParam string

func (p URLParam) String() string { return string(p) }

func (p URLParam) Validate() error {
	if strings.TrimSpace(string(p)) == "" {
		return fmt.Errorf("param %s is empty", p)
	}
	return nil
}
