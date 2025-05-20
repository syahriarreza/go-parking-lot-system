
package utils

import "fmt"

// FormatSpot returns string like "M-1(X)" or "A-1(O)"
func FormatSpot(spotType string, occupied bool) string {
	if occupied {
		return fmt.Sprintf("%s(X)", spotType)
	}
	return fmt.Sprintf("%s(O)", spotType)
}
