package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Parse(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %v", path, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("formato non valido alla riga %d: %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Rimuovi le virgolette se presenti
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		err = os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("errore nell'impostare la variabile %s: %v", key, err)
		}
	}

	return nil
}
