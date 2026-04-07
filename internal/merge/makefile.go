package merge

import (
	"bufio"
	"bytes"
	"regexp"
)

var targetRegex = regexp.MustCompile(`^([a-zA-Z0-9_.-]+):`)

func getTargets(content []byte) []string {
	var targets []string
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := targetRegex.FindStringSubmatch(line); len(matches) > 1 {
			// Don't include .PHONY
			if matches[1] != ".PHONY" {
				targets = append(targets, matches[1])
			}
		}
	}
	return targets
}

func getTargetBlocks(content []byte) map[string]string {
	blocks := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(content))

	var currentTarget string
	var currentBlock bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		matches := targetRegex.FindStringSubmatch(line)

		if len(matches) > 1 && matches[1] != ".PHONY" {
			if currentTarget != "" {
				blocks[currentTarget] = currentBlock.String()
				currentBlock.Reset()
			}
			currentTarget = matches[1]
			currentBlock.WriteString(line + "\n")
		} else if currentTarget != "" {
			currentBlock.WriteString(line + "\n")
		}
	}

	if currentTarget != "" {
		blocks[currentTarget] = currentBlock.String()
	}

	return blocks
}

// Makefile safely appends missing targets from src into dst.
func Makefile(dst, src []byte) ([]byte, error) {
	dstTargets := getTargets(dst)
	dstSet := make(map[string]bool)
	for _, t := range dstTargets {
		dstSet[t] = true
	}

	srcBlocks := getTargetBlocks(src)
	srcTargets := getTargets(src)

	var newBlocks bytes.Buffer
	for _, t := range srcTargets {
		if !dstSet[t] {
			newBlocks.WriteString("\n" + srcBlocks[t])
		}
	}

	if newBlocks.Len() > 0 {
		return append(dst, newBlocks.Bytes()...), nil
	}

	return dst, nil
}
