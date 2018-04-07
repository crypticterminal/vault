package apidoc

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	var ess []string

	tests := []struct {
		root         string
		in_pattern   string
		out_pathlets []pathlet
	}{
		{"sys", "rekey/backup", []pathlet{{"/sys/rekey/backup", ess}}},
		{"sys", "rekey/backup$", []pathlet{{"/sys/rekey/backup", ess}}},
		{"sys", "auth/(?P<path>.+?)/tune$", []pathlet{{"/sys/auth/{path}/tune", set("path")}}},
		{"sys", "auth/(?P<path>.+?)/tune/(?P<more>.*?)$", []pathlet{{"/sys/auth/{path}/tune/{more}", set("path", "more")}}},
		{"sys", "tools/hash(/(?P<urlalgorithm>.+))?", []pathlet{
			{"/sys/tools/hash", set()},
			{"/sys/tools/hash/{urlalgorithm}", set("urlalgorithm")},
		}},
		{"sys", "(leases/)?renew(/(?P<url_lease_id>.+))?", []pathlet{
			{"/sys/leases/renew", set()},
			{"/sys/leases/renew/{url_lease_id}", set("url_lease_id")},
			{"/sys/renew", set()},
			{"/sys/renew/{url_lease_id}", set("url_lease_id")},
		}},
		{"sys", `config/ui/headers/(?P<header>\w(([\w-.]+)?\w)?)`, []pathlet{
			{"/sys/config/ui/headers/{header}", set("header")},
		}},

		/* optional elements
		tools/hash(/(?P<urlalgorithm>.+))?
		{"sys", "leases/lookup/(?P<prefix>.+?)?", []pathlet{
			{"/sys/leases/lookup", []string{"path", "more"}},
			{"/sys/leases/lookup/{prefix}", []string{"prefix"}},
		}},
		*/
		//"leases/lookup/(?P<prefix>.+?)?"
		//"(leases/)?renew"
	}
	/*
		root := "sys"
		pat := "rekey/backup"
		exp := "/sys/rekey/backup"
		out := parsePattern(root, pat)[0].pattern
	*/
	for i, test := range tests {
		out := expandPattern(test.root, test.in_pattern)
		if !reflect.DeepEqual(out, test.out_pathlets) {
			t.Fatalf("Test %d: Expected %v got %v", i, test.out_pathlets, out)
		}
	}
}

func set(strings ...string) []string {
	return strings
}
