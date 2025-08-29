package wt

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type WorktreeItem struct {
	Path   string
	Head   string
	Branch string // short name, no refs/heads/
}

func (w WorktreeItem) String() string {
	return fmt.Sprintf("WorktreeItem{Path: %s, Branch: %s, Head: %s}", w.Path, w.Branch, w.Head)
}

func List() ([]WorktreeItem, error) {
	if _, err := repoRoot(); err != nil {
		return nil, err
	}

	out, err := gitOut("worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}

	var items []WorktreeItem
	var cur WorktreeItem

	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "worktree "):
			if cur.Path != "" {
				items = append(items, cur)
			}
			cur = WorktreeItem{Path: strings.TrimSpace(strings.TrimPrefix(line, "worktree "))}
		case strings.HasPrefix(line, "HEAD "):
			cur.Head = strings.TrimSpace(strings.TrimPrefix(line, "HEAD "))
		case strings.HasPrefix(line, "branch "):
			ref := strings.TrimSpace(strings.TrimPrefix(line, "branch "))
			cur.Branch = strings.TrimPrefix(ref, "refs/heads/")
		}
	}
	if cur.Path != "" {
		items = append(items, cur)
	}
	return items, nil
}
