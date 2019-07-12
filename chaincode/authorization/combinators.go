package authorization

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/context"
)

func All(cc ...Check) Check {
	return func(ctx *context.Context) error {
		for _, c := range cc {
			if err := c(ctx); err != nil {
				return errors.Wrap(err, "All: one did not allow")
			}
		}
		return nil
	}
}

func Any(cc ...Check) Check {
	return func(ctx *context.Context) error {
		errs := []string{}
		for _, c := range cc {
			err := c(ctx)
			if err == nil {
				return nil
			}
			errs = append(errs, err.Error())
		}
		return errors.Errorf("Any: all did not allow (" + strings.Join(errs, ",") + ")")
	}
}
