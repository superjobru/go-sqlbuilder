// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import (
	"fmt"
	"testing"

	"github.com/huandu/go-assert"
)

func TestFlavor(t *testing.T) {
	a := assert.New(t)
	cases := map[Flavor]string{
		0:        "<invalid>",
		SphinxQL: "SphinxQL",
	}

	for f, expected := range cases {
		actual := f.String()
		a.Equal(actual, expected)
	}
}

func ExampleFlavor_Interpolate_sphinxQL() {
	sb := SphinxQL.NewSelectBuilder()
	sb.Select("name").From("user").Where(
		sb.NE("id", 1234),
		sb.E("name", "Charmy Liu"),
		sb.Like("desc", "%mother's day%"),
	)
	sql, args := sb.Build()
	query, err := SphinxQL.Interpolate(sql, args)

	fmt.Println(query)
	fmt.Println(err)

	// Output:
	// SELECT name FROM user WHERE id <> 1234 AND name = 'Charmy Liu' AND desc LIKE '%mother\'s day%'
	// <nil>
}
