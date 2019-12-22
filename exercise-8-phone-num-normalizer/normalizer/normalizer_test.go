package normalizer

import (
	"testing"
)

func TestSpaceRemoval(t *testing.T) {
	testInput := "1234  5   "
	expectedOutput := "12345"
	checkResults(testInput, expectedOutput, t)
}

func TestDashRemoval(t *testing.T) {
	testInput := "123-4-5"
	expectedOutput := "12345"
	checkResults(testInput, expectedOutput, t)
}
func TestBracketRemoval(t *testing.T) {
	testInput := "123(4)5"
	expectedOutput := "12345"
	checkResults(testInput, expectedOutput, t)
}

func TestItAll(t *testing.T) {
	multipleInputs := map[string]string{
		"1234567890":     "1234567890",
		"123 456 7891":   "1234567891",
		"(123) 456 7892": "1234567892",
		"(123) 456-7893": "1234567893",
		"123-456-7894":   "1234567894",
		"123-456-7890":   "1234567890",
		"1234567892":     "1234567892",
		"(123)456-7892":  "1234567892",
	}
	for testInput, expectedOutput := range multipleInputs {
		checkResults(testInput, expectedOutput, t)
	}
}

func checkResults(inp string, out string, t *testing.T) {
	normalized, err := NormalizeNumber(inp)
	if err != nil {
		t.Fatalf("Failed: %s is not a number-there are other special characters not normally found in phone numbers", inp)
	}
	if normalized == out {
		t.Logf("Success: %s normalized to %s", inp, out)
	} else {
		t.Errorf("Failed: %s normalized to %s instead of %s", inp, normalized, out)
	}
}
