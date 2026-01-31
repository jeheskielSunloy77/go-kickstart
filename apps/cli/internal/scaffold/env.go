package scaffold

import (
	"bufio"
	"strings"
)

func MergeEnvExample(input string, overrides map[string]string) string {
	seen := map[string]bool{}
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") || !strings.Contains(line, "=") {
			out = append(out, line)
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := parts[1]
		if override, ok := overrides[key]; ok {
			val = override
			seen[key] = true
		}
		out = append(out, key+"="+val)
	}
	for key, val := range overrides {
		if !seen[key] {
			out = append(out, key+"="+val)
		}
	}
	return strings.Join(out, "\n") + "\n"
}
