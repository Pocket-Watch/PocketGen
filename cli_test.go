package main

import "testing"

func TestChangeExtension(t *testing.T) {
	actual := changeExtension("example.txt", ".go")
	expected := "example.go"
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestChangeExtensionNoExtension(t *testing.T) {
	actual := changeExtension("example", ".rs")
	expected := "example.rs"
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestChangeExtensionTwoDots(t *testing.T) {
	actual := changeExtension("example.txt.tg", ".rs")
	expected := "example.txt.rs"
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
