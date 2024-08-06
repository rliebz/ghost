package be

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"

	"github.com/rliebz/ghost/internal/color"
	"github.com/rliebz/ghost/internal/jsondiff"
)

var exportTypes = cmp.Exporter(func(reflect.Type) bool { return true })

func colorDiff[T any](x, y T, opts ...cmp.Option) string {
	diff := cmp.Diff(x, y, append(opts, exportTypes)...)
	return applyColors(diff)
}

func colorJSONDiff[T ~string | ~[]byte](got, want T) (string, jsondiff.Kind) {
	diff, kind := jsondiff.Diff(got, want)
	return applyColors(diff), kind
}

func applyColors(diff string) string {
	if diff == "" {
		return ""
	}

	ss := strings.Split(diff, "\n")
	for i, s := range ss {
		switch {
		case strings.HasPrefix(s, "-"):
			ss[i] = color.Red(s)
		case strings.HasPrefix(s, "+"):
			ss[i] = color.Green(s)
		// Only color the first character, since we expect inline red/green
		case strings.HasPrefix(s, "~"):
			ss[i] = color.Yellow("~") + s[1:]
		}
	}

	return fmt.Sprintf(
		`diff (%s %s):
%v`,
		color.Red("-want"),
		color.Green("+got"),
		strings.Join(ss, "\n"),
	)
}
