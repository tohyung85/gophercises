package normalizer

import (
	"strconv"
	"strings"
)

func NormalizeNumber(num string) (string, error) {
	unwantedChars := [4]string{" ", "(", ")", "-"}

	for _, c := range unwantedChars {
		num = strings.ReplaceAll(num, c, "")
	}

	_, err := strconv.Atoi(num)
	if err != nil {
		return "", err
	}

	return num, nil
}
