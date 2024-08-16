package naming

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

var validSymbols = []rune{
	'_', '-',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

const symbolForReplacement = '_'

func GenContainerName(imageTag string, workDir string) string {
	return fmt.Sprintf("%s_%s",
		GenImageName(imageTag),
		filepath.Base(workDir),
	)
}

func GenImageName(imageTag string) string {
	b := strings.Builder{}
	for _, ch := range strings.ToLower(imageTag) {
		if slices.Contains(validSymbols, ch) {
			b.WriteRune(ch)
			continue
		}
		b.WriteRune(symbolForReplacement)
	}
	return b.String()
}

func GenDevHomeDirName(imageTag string, workDir string) string {
	return fmt.Sprintf(".%s--%s--dev-home",
		filepath.Base(workDir),
		GenImageName(imageTag),
	)
}
