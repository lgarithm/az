package cloudinit

import (
	"errors"
	"log"
	"strings"
)

// handlers.py[DEBUG]: finish: modules-final/config-power-state-change: SUCCESS: config-power-state-change ran successfully
// handlers.py[DEBUG]: finish: modules-final: SUCCESS: running modules for final
const lastLineMark = `handlers.py[DEBUG]: finish: modules-final: `
const successLastLineMark = `handlers.py[DEBUG]: finish: modules-final: SUCCESS`

func IsLastLine(line string) bool {
	return strings.Contains(line, lastLineMark)
}

func FinalError(last string) error {
	if strings.Contains(last, successLastLineMark) {
		return nil
	}
	parts := strings.Split(last, lastLineMark)
	if len(parts) > 1 {
		return errors.New(parts[1])
	}
	log.Printf("unexpected cloud-init log format: %q", last)
	return errors.New(last)
}
