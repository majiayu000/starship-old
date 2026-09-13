package postgres

import (
	"errors"
	"testing"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

func TestResolveSortColumn_AllowlistedColumns(t *testing.T) {
	cases := []struct {
		name       string
		dataSource string
		sortBy     string
		want       string
	}{
		{"camelCase domain", domain.DataSourceSATOneprep, "domain", "domain"},
		{"camelCase originalId", domain.DataSourceSATOneprep, "originalId", "original_id"},
		{"snake_case original_id", domain.DataSourceSATOneprep, "original_id", "original_id"},
		{"camelCase reviewStatus", domain.DataSourceSATOneprep, "reviewStatus", "review_status"},
		{"snake_case review_status", domain.DataSourceSATOneprep, "review_status", "review_status"},
		{"ixl skill", domain.DataSourceSATIXL, "skill", "skill"},
		{"ixl knowledgePoint", domain.DataSourceSATIXL, "knowledgePoint", "knowledge_point"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.ResolveSortColumn(tc.dataSource, tc.sortBy)
			if err != nil {
				t.Fatalf("ResolveSortColumn(%q, %q) unexpected error: %v", tc.dataSource, tc.sortBy, err)
			}
			if got != tc.want {
				t.Fatalf("ResolveSortColumn(%q, %q) = %q, want %q", tc.dataSource, tc.sortBy, got, tc.want)
			}
		})
	}
}

func TestResolveSortColumn_DefaultWhenEmpty(t *testing.T) {
	got, err := domain.ResolveSortColumn(domain.DataSourceSATOneprep, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "original_id" {
		t.Fatalf("empty sortBy = %q, want original_id", got)
	}
}

func TestResolveSortColumn_RejectsInjectionPayloads(t *testing.T) {
	payloads := []string{
		"1;DROP TABLE users",
		"(SELECT 1)",
		"domain ASC; DELETE FROM t_math_sat_temp_classify_oneprep --",
		"domain, (SELECT pg_sleep(1))",
		"original_id; SELECT 1",
		" domain ",
		"domain ",
		" domain",
		"unknownColumn",
		"review_comment", // not in allowlist for ORDER BY exposure
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			got, err := domain.ResolveSortColumn(domain.DataSourceSATOneprep, payload)
			if err == nil {
				t.Fatalf("expected error for sortBy %q, got column %q", payload, got)
			}
			if !errors.Is(err, domain.ErrInvalidSortBy) {
				t.Fatalf("expected ErrInvalidSortBy for %q, got %v", payload, err)
			}
			if got != "" {
				t.Fatalf("expected empty column on error, got %q", got)
			}
		})
	}
}
