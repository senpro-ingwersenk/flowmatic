package workgroup_test

import (
	"errors"
	"testing"

	"github.com/carlmjohnson/workgroup"
)

func TestDoAll_err(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")
	errs := workgroup.DoAll(1, []int{1, 2, 3}, func(i int) error {
		switch i {
		case 1:
			return a
		case 2:
			return b
		default:
			return nil
		}
	})
	if !errors.Is(errs, a) {
		t.Fatal(errs)
	}
	if !errors.Is(errs, b) {
		t.Fatal(errs)
	}
}
