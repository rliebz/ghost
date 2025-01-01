package jsondiff_test

import (
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
	"github.com/rliebz/ghost/internal/jsondiff"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantDiff string
		wantKind jsondiff.Kind
	}{
		// Primitives
		{
			name:     "nulls equal",
			a:        `null`,
			b:        `null`,
			wantDiff: `null`,
			wantKind: jsondiff.Match,
		},
		{
			name:     "bools equal",
			a:        `true`,
			b:        `true`,
			wantDiff: `true`,
			wantKind: jsondiff.Match,
		},
		{
			name:     "bools not equal",
			a:        `true`,
			b:        `false`,
			wantDiff: `false => true`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name:     "string equal",
			a:        `"foo"`,
			b:        `"foo"`,
			wantDiff: `"foo"`,
			wantKind: jsondiff.Match,
		},
		{
			name:     "string not equal",
			a:        `"foo"`,
			b:        `"bar"`,
			wantDiff: `"bar" => "foo"`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name:     "string and null",
			a:        `"null"`,
			b:        `null`,
			wantDiff: `null => "null"`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name:     "number equal",
			a:        `17.5`,
			b:        `17.5`,
			wantDiff: `17.5`,
			wantKind: jsondiff.Match,
		},
		{
			name:     "number not equal",
			a:        `17.5`,
			b:        `17.6`,
			wantDiff: `17.6 => 17.5`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name:     "number and string",
			a:        `17.5`,
			b:        `"17.5"`,
			wantDiff: `"17.5" => 17.5`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name:     "got invalid",
			a:        `NaN`,
			b:        `true`,
			wantKind: jsondiff.GotInvalid,
		},
		{
			name:     "want invalid",
			a:        `true`,
			b:        `NaN`,
			wantKind: jsondiff.WantInvalid,
		},
		{
			name:     "both invalid",
			a:        `NaN`,
			b:        `NaN`,
			wantKind: jsondiff.BothInvalid,
		},

		// Objects
		{
			name:     "empty objects equal",
			a:        `{}`,
			b:        `{}`,
			wantDiff: `{}`,
			wantKind: jsondiff.Match,
		},
		{
			name: "small objects equal",
			a:    `{"a": 1}`,
			b:    `{"a": 1}`,
			wantDiff: `  {
    "a": 1
  }`,
			wantKind: jsondiff.Match,
		},
		{
			name: "small objects values not equal",
			a:    `{"a": 1}`,
			b:    `{"a": "1"}`,
			wantDiff: `  {
~   "a": "1" => 1
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "small objects keys not equal",
			a:    `{"a": 1}`,
			b:    `{"b": 1}`,
			wantDiff: `  {
+   "a": 1,
-   "b": 1
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "small objects extra elements",
			a:    `{"a": 1}`,
			b:    `{}`,
			wantDiff: `  {
+   "a": 1
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "small objects missing elements",
			a:    `{}`,
			b:    `{"b": 1}`,
			wantDiff: `  {
-   "b": 1
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "large objects equal unordered",
			a:    `{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`,
			b:    `{"c": 3, "b": 2, "a": 1, "e": 5, "d": 4}`,
			wantDiff: `  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": 4,
    "e": 5
  }`,
			wantKind: jsondiff.Match,
		},

		// Arrays
		{
			name:     "empty arrays equal",
			a:        `[]`,
			b:        `[]`,
			wantDiff: `[]`,
			wantKind: jsondiff.Match,
		},
		{
			name: "small arrays equal",
			a:    `["a"]`,
			b:    `["a"]`,
			wantDiff: `  [
    "a"
  ]`,
			wantKind: jsondiff.Match,
		},
		{
			name: "small arrays not equal",
			a:    `["a"]`,
			b:    `["b"]`,
			wantDiff: `  [
~   "b" => "a"
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "small arrays missing element",
			a:    `[]`,
			b:    `["b"]`,
			wantDiff: `  [
-   "b"
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "small arrays extra element",
			a:    `["a"]`,
			b:    `[]`,
			wantDiff: `  [
+   "a"
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "arrays respect order",
			a:    `["a", "b"]`,
			b:    `["b", "a"]`,
			wantDiff: `  [
~   "b" => "a",
~   "a" => "b"
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "large arrays equal",
			a:    `["a", "b", "c", "d", "e"]`,
			b:    `["a", "b", "c", "d", "e"]`,
			wantDiff: `  [
    "a",
    "b",
    "c",
    "d",
    "e"
  ]`,
			wantKind: jsondiff.Match,
		},

		// Compound structures
		{
			name:     "empty object not equal empty array",
			a:        `{}`,
			b:        `[]`,
			wantDiff: `[] => {}`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name:     "object not equal array",
			a:        `{"a": 1}`,
			b:        `["a", 1]`,
			wantDiff: `[...] => {...}`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "objects containing array equal",
			a:    `{"a": [1, 2, 3]}`,
			b:    `{"a": [1, 2, 3]}`,
			wantDiff: `  {
    "a": [
      1,
      2,
      3
    ]
  }`,
			wantKind: jsondiff.Match,
		},
		{
			name: "objects containing array not equal",
			a:    `{"a": [1, 3, 2]}`,
			b:    `{"a": [1, 2, 3]}`,
			wantDiff: `  {
    "a": [
      1,
~     2 => 3,
~     3 => 2
    ]
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "object containing empty array",
			a:    `{"a": []}`,
			b:    `{}`,
			wantDiff: `  {
+   "a": []
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "object missing empty array",
			a:    `{}`,
			b:    `{"a": []}`,
			wantDiff: `  {
-   "a": []
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "object containing extra array",
			a:    `{"a": [1, 2, 3]}`,
			b:    `{}`,
			wantDiff: `  {
+   "a": [
+     1,
+     2,
+     3
+   ]
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "object missing extra array",
			a:    `{}`,
			b:    `{"a": [1, 2, 3]}`,
			wantDiff: `  {
-   "a": [
-     1,
-     2,
-     3
-   ]
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "object containing array containing object equal",
			a:    `{"a": [1, {"b": 2}, 3]}`,
			b:    `{"a": [1, {"b": 2}, 3]}`,
			wantDiff: `  {
    "a": [
      1,
      {
        "b": 2
      },
      3
    ]
  }`,
			wantKind: jsondiff.Match,
		},
		{
			name: "object containing array containing object not equal",
			a:    `{"a": [1, {"b": 2}, 3]}`,
			b:    `{"a": [1, {"b": 4}, 3]}`,
			wantDiff: `  {
    "a": [
      1,
      {
~       "b": 4 => 2
      },
      3
    ]
  }`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "array containing object containing array equal",
			a:    `[{"a": [1, 2, 3]}]`,
			b:    `[{"a": [1, 2, 3]}]`,
			wantDiff: `  [
    {
      "a": [
        1,
        2,
        3
      ]
    }
  ]`,
			wantKind: jsondiff.Match,
		},
		{
			name: "array containing object containing array not equal",
			a:    `[{"a": [1, 2, 3]}]`,
			b:    `[{"a": [1, 3, 2]}]`,
			wantDiff: `  [
    {
      "a": [
        1,
~       3 => 2,
~       2 => 3
      ]
    }
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "array containing empty object",
			a:    `[0, {}]`,
			b:    `[0]`,
			wantDiff: `  [
    0,
+   {}
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "array missing empty object",
			a:    `[0]`,
			b:    `[0, {}]`,
			wantDiff: `  [
    0,
-   {}
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "array containing extra object",
			a:    `[0, {"a": 1}]`,
			b:    `[0]`,
			wantDiff: `  [
    0,
+   {
+     "a": 1
+   }
  ]`,
			wantKind: jsondiff.NoMatch,
		},
		{
			name: "array missing extra object",
			a:    `[0]`,
			b:    `[0, {"a": 1}]`,
			wantDiff: `  [
    0,
-   {
-     "a": 1
-   }
  ]`,
			wantKind: jsondiff.NoMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := ghost.New(t)

			diff, kind := jsondiff.Diff(tt.a, tt.b)
			g.Should(be.Equal(kind, tt.wantKind))
			g.Should(be.Equal(diff, tt.wantDiff))
		})
	}
}

func TestKind_String(t *testing.T) {
	g := ghost.New(t)

	k := jsondiff.Kind(-1)
	g.Should(be.Equal(k.String(), "InvalidKind"))
}
