package api_test

import (
	"testing"

	"github.com/mattribution/api/pkg/api"
)

func TestInvalidKPI(t *testing.T) {
	kpi := api.KPI{}
	expected := false
	got := kpi.IsValid()
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestValidKPI(t *testing.T) {
	kpi := api.KPI{
		Column: "test",
		Value:  "test",
	}
	expected := false
	got := kpi.IsValid()
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
