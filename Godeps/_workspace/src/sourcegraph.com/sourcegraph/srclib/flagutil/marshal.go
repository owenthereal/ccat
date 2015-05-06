package flagutil

import (
	"fmt"
	"strings"

	"sourcegraph.com/sourcegraph/go-flags"
)

// MarshalArgs takes a struct with go-flags field tags and turns it into an
// equivalent []string for use as command-line args.
func MarshalArgs(v interface{}) ([]string, error) {
	parser := flags.NewParser(nil, flags.None)
	group, err := parser.AddGroup("", "", v)
	if err != nil {
		return nil, err
	}
	return marshalArgsInGroup(group, "")
}

func marshalArgsInGroup(group *flags.Group, prefix string) ([]string, error) {
	var args []string
	for _, opt := range group.Options() {
		flagStr := opt.String()

		// handle flags with both short and long (just get the long)
		if i := strings.Index(flagStr, ", --"); i != -1 {
			flagStr = flagStr[i+2:]
		}

		v := opt.Value()
		if m, ok := v.(flags.Marshaler); ok {
			s, err := m.MarshalFlag()
			if err != nil {
				return nil, err
			}
			args = append(args, flagStr, s)
		} else if ss, ok := v.([]string); ok {
			for _, s := range ss {
				args = append(args, flagStr, s)
			}
		} else if bv, ok := v.(bool); ok {
			if bv {
				args = append(args, flagStr)
			}
		} else {
			args = append(args, flagStr, fmt.Sprintf("%v", opt.Value()))
		}
	}
	for _, g := range group.Groups() {
		// TODO(sqs): assumes that the NamespaceDelimiter is "."
		const namespaceDelimiter = "."
		groupArgs, err := marshalArgsInGroup(g, g.Namespace+namespaceDelimiter)
		if err != nil {
			return nil, err
		}
		args = append(args, groupArgs...)
	}
	return args, nil
}
