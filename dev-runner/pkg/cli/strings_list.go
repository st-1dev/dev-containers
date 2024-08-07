package cli

import (
	"fmt"
	"sort"
)

type StringSlice struct {
	sort.StringSlice
}

func (s *StringSlice) String() string {
	return fmt.Sprintf("%#v", s.StringSlice)
}

func (s *StringSlice) Set(value string) error {
	s.StringSlice = append(s.StringSlice, value)
	return nil
}
