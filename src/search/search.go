package search

import (
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/rostislavjadavan/mdwiki/src/storage"
	"github.com/sahilm/fuzzy"
	stripmd "github.com/writeas/go-strip-markdown"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Result struct {
	Query       string         `json:"query"`
	Filenames   []searchResult `json:"filenames"`
	PageContent []searchResult `json:"pageContent"`
}

type searchResult struct {
	Filename string    `json:"filename"`
	ModTime  time.Time `json:"modTime"`
	Score    int       `json:"score"`
	Preview  string    `json:"preview"`
}

const maxPreviewLines = 3

func Search(query string, s *storage.Storage) (*Result, error) {
	result := Result{
		Query:       query,
		Filenames:   make([]searchResult, 0),
		PageContent: make([]searchResult, 0),
	}

	pages, err := s.PageList()
	if err != nil {
		return &result, err
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return &result, nil
	}

	normalizedQuery := normalizeString(query)

	// Filename matching
	filenames := make([]string, len(pages))
	for i, p := range pages {
		filenames[i] = normalizeString(p.Filename)
	}

	matches := fuzzy.Find(normalizedQuery, filenames)
	for _, m := range matches {
		page := pages[m.Index]
		preview := highlightFromPositions(page.Filename, m.MatchedIndexes)
		result.Filenames = append(result.Filenames, searchResult{
			Filename: page.Filename,
			ModTime:  page.ModTime,
			Score:    m.Score,
			Preview:  preview,
		})
	}

	// Content matching
	for _, page := range pages {
		markdownContent, _ := s.PageRawContent(page.Filename)
		content := stripmd.Strip(markdownContent)
		lines := strings.Split(content, "\n")

		normalizedLines := make([]string, 0, len(lines))
		originalLines := make([]string, 0, len(lines))
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			normalizedLines = append(normalizedLines, normalizeString(trimmed))
			originalLines = append(originalLines, trimmed)
		}

		lineMatches := fuzzy.Find(normalizedQuery, normalizedLines)

		if len(lineMatches) == 0 {
			continue
		}

		r := searchResult{
			Filename: page.Filename,
			ModTime:  page.ModTime,
			Score:    0,
			Preview:  "",
		}

		count := 0
		for _, lm := range lineMatches {
			if count >= maxPreviewLines {
				break
			}
			r.Score += lm.Score
			preview := highlightFromPositions(originalLines[lm.Index], lm.MatchedIndexes)
			if r.Preview != "" {
				r.Preview += "\n"
			}
			r.Preview += preview
			count++
		}

		result.PageContent = append(result.PageContent, r)
	}

	sort.Slice(result.PageContent, func(i, j int) bool {
		return result.PageContent[i].Score > result.PageContent[j].Score
	})

	sort.Slice(result.Filenames, func(i, j int) bool {
		return result.Filenames[i].Score > result.Filenames[j].Score
	})

	return &result, nil
}

func normalizeString(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, strings.ToLower(s))
	return result
}

func highlightFromPositions(s string, positions []int) string {
	if len(positions) == 0 {
		return s
	}

	posSet := make(map[int]bool, len(positions))
	for _, p := range positions {
		posSet[p] = true
	}

	var b strings.Builder
	inHighlight := false
	for i, ch := range []rune(s) {
		if posSet[i] {
			if !inHighlight {
				b.WriteString(`<b class="highlight">`)
				inHighlight = true
			}
			b.WriteRune(ch)
		} else {
			if inHighlight {
				b.WriteString("</b>")
				inHighlight = false
			}
			b.WriteRune(ch)
		}
	}
	if inHighlight {
		b.WriteString("</b>")
	}
	return b.String()
}
