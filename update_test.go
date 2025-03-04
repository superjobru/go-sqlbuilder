// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
	"testing"

	"github.com/huandu/go-assert"
)

func ExampleUpdate() {
	sql := Update("demo.user").
		Set(
			"visited = visited + 1",
		).
		Where(
			"id = 1234",
		).
		String()

	fmt.Println(sql)

	// Output:
	// UPDATE demo.user SET visited = visited + 1 WHERE id = 1234
}

func ExampleUpdateBuilder() {
	ub := NewUpdateBuilder()
	ub.Update("demo.user")
	ub.Set(
		ub.Assign("type", "sys"),
		ub.Incr("credit"),
		"modified_at = UNIX_TIMESTAMP(NOW())", // It's allowed to write arbitrary SQL.
	)
	ub.Where(
		ub.GreaterThan("id", 1234),
		ub.Like("name", "%Du"),
		ub.Or(
			ub.IsNull("id_card"),
			ub.In("status", 1, 2, 5),
		),
		"modified_at > created_at + "+ub.Var(86400), // It's allowed to write arbitrary SQL.
	)
	ub.Option(
		ub.Comment("kekw"),
	)

	sql, args := ub.Build()
	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// UPDATE demo.user SET type = ?, credit = credit + 1, modified_at = UNIX_TIMESTAMP(NOW()) WHERE id > ? AND name LIKE ? AND (id_card IS NULL OR status IN (?, ?, ?)) AND modified_at > created_at + ? OPTION comment = ?
	// [sys 1234 %Du 1 2 5 86400 kekw]
}

func TestUpdateAssignments(t *testing.T) {
	a := assert.New(t)
	cases := map[string]func(ub *UpdateBuilder) string{
		"f = f + 1|[]":     func(ub *UpdateBuilder) string { return ub.Incr("f") },
		"f = f - 1|[]":     func(ub *UpdateBuilder) string { return ub.Decr("f") },
		"f = f + $0|[123]": func(ub *UpdateBuilder) string { return ub.Add("f", 123) },
		"f = f - $0|[123]": func(ub *UpdateBuilder) string { return ub.Sub("f", 123) },
		"f = f * $0|[123]": func(ub *UpdateBuilder) string { return ub.Mul("f", 123) },
		"f = f / $0|[123]": func(ub *UpdateBuilder) string { return ub.Div("f", 123) },
	}

	for expected, f := range cases {
		ub := NewUpdateBuilder()
		s := f(ub)
		ub.Set(s)
		_, args := ub.Build()
		actual := fmt.Sprintf("%v|%v", s, args)

		a.Equal(actual, expected)
	}
}

func ExampleUpdateBuilder_SetMore() {
	ub := NewUpdateBuilder()
	ub.Update("demo.user")
	ub.Set(
		ub.Assign("type", "sys"),
		ub.Incr("credit"),
	)
	ub.SetMore(
		"modified_at = UNIX_TIMESTAMP(NOW())", // It's allowed to write arbitrary SQL.
	)

	sql, args := ub.Build()
	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// UPDATE demo.user SET type = ?, credit = credit + 1, modified_at = UNIX_TIMESTAMP(NOW())
	// [sys]
}

func ExampleUpdateBuilder_SQL() {
	ub := NewUpdateBuilder()
	ub.SQL("/* before */")
	ub.Update("demo.user")
	ub.SQL("/* after update */")
	ub.Set(
		ub.Assign("type", "sys"),
	)
	ub.SQL("/* after set */")

	sql := ub.String()
	fmt.Println(sql)

	// Output:
	// /* before */ UPDATE demo.user /* after update */ SET type = ? /* after set */
}
