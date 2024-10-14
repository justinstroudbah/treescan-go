package treescanner

import "testing"

func TestScan(t *testing.T) {
	scanManager := NewScanManager()

	if scanManager == nil {
		t.Fatal("scanManager is nil")
	}

}

func TestLoadConfig(t *testing.T) {

}
