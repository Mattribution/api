package api_test

import (
	"testing"

	"github.com/mattribution/api/pkg/api"
)

func TestAdjustWeightAddOne(t *testing.T) {
	kpi := api.KPI{}
	kpi.AdjustWeight("firstTouch", "mock-key", 1)
	expectedDataStr := `{"firstTouch":{"weights":{"mock-key":1}}}`
	gotStr := kpi.Data.String()
	if gotStr != expectedDataStr {
		t.Errorf("Expected %v, got %v", expectedDataStr, gotStr)
	}
}

func TestAdjustWeightAddTwice(t *testing.T) {
	kpi := api.KPI{}
	kpi.AdjustWeight("firstTouch", "mock-key", 1)
	kpi.AdjustWeight("firstTouch", "mock-key", 1)
	expectedDataStr := `{"firstTouch":{"weights":{"mock-key":2}}}`
	gotStr := kpi.Data.String()
	if gotStr != expectedDataStr {
		t.Errorf("Expected %v, got %v", expectedDataStr, gotStr)
	}
}

func TestAdjustWeightMultipleKeys(t *testing.T) {
	kpi := api.KPI{}
	kpi.AdjustWeight("firstTouch", "mock-key", 1)
	kpi.AdjustWeight("firstTouch", "mock-key-2", 1)
	expectedDataStr := `{"firstTouch":{"weights":{"mock-key":1,"mock-key-2":1}}}`
	gotStr := kpi.Data.String()
	if gotStr != expectedDataStr {
		t.Errorf("Expected %v, got %v", expectedDataStr, gotStr)
	}
}

func TestAdjustWeightNegative(t *testing.T) {
	kpi := api.KPI{}
	kpi.AdjustWeight("firstTouch", "mock-key", 1)
	kpi.AdjustWeight("firstTouch", "mock-key", -1)
	expectedDataStr := `{"firstTouch":{"weights":{"mock-key":0}}}`
	gotStr := kpi.Data.String()
	if gotStr != expectedDataStr {
		t.Errorf("Expected %v, got %v", expectedDataStr, gotStr)
	}
}

func TestClearWeightsForModel(t *testing.T) {
	kpi := api.KPI{}
	kpi.AdjustWeight("firstTouch", "mock-key", 1)
	kpi.ClearWeightsForModel("firstTouch")
	expectedDataStr := `{}`
	gotStr := kpi.Data.String()
	if gotStr != expectedDataStr {
		t.Errorf("Expected %v, got %v", expectedDataStr, gotStr)
	}
}
