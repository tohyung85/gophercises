package normalizer

import "regexp"

// func NormalizeNumber(num string) (string, error) {
// 	unwantedChars := [4]string{" ", "(", ")", "-"}

// 	for _, c := range unwantedChars {
// 		num = strings.ReplaceAll(num, c, "")
// 	}

// 	_, err := strconv.Atoi(num)
// 	if err != nil {
// 		return "", err
// 	}

// 	return num, nil
// }

// func NormalizeNumber(num string) (string, error) { // More robust
// 	var buf bytes.Buffer

// 	for _, ch := range num {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String(), nil
// }

func NormalizeNumber(num string) (string, error) { // Using Regex
	re := regexp.MustCompile("[^0-9]")
	normalized := re.ReplaceAllString(num, "")
	return normalized, nil
}
