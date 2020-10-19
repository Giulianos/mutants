package util

import "testing"

func TestMaxInt(t *testing.T) {
	tests := []struct{
		n1, n2, expected int
	}{
		{1,2,2},
		{2,1,2},
		{1,1,1},
	}

	for _, test := range tests {
		if MaxInt(test.n1, test.n2) != test.expected {
			t.Errorf("MaxInt(%d, %d) must be %d", test.n1, test.n2, test.expected)
		}
	}
}

func TestMinInt(t *testing.T) {
	tests := []struct{
		n1, n2, expected int
	}{
		{1,2,1},
		{2,1,1},
		{1,1,1},
	}

	for _, test := range tests {
		if MinInt(test.n1, test.n2) != test.expected {
			t.Errorf("MinInt(%d, %d) must be %d", test.n1, test.n2, test.expected)
		}
	}
}
