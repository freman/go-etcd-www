package main

import (
	"fmt"
	"strings"
)

type flagStrs struct{
	Values []string
	reset bool
}

func newFlagStrs(defaults []string) *flagStrs {
	return &flagStrs{defaults, true};
}

func (s *flagStrs) String() string {
	return fmt.Sprint(s.Values)
}

func (s *flagStrs) Set(value string) error {
	if strings.Contains(value, ",") {
		s.Values = strings.Split(value, ",")
	} else {
		if (s.reset) {
			s.Values = make([]string, 0)
			s.reset = false
		}
		s.Values = append(s.Values, value)
	}
	return nil
}
