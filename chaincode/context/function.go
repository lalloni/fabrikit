package context

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func ParseFunction(bs []byte) (string, map[string]string, error) {
	ss := strings.SplitN(string(bs), "?", 2)
	fun := ss[0]
	opt := map[string]string{}
	if len(ss) == 1 {
		return fun, opt, nil
	}
	values, err := url.ParseQuery(ss[1])
	if err != nil {
		return fun, opt, errors.Wrap(err, "parsing function options")
	}
	for n := range values {
		opt[n] = values.Get(n)
	}
	return fun, opt, nil
}
