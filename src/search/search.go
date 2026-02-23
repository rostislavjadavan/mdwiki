package search

import (
	"github.com/rostislavjadavan/mdwiki/src/storage"
	"github.com/writeas/go-strip-markdown"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"sort"
	"strings"
	"time"
	"unicode"
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

	queries := extractQueries(query)
	if len(queries) == 0 {
		return &result, nil
	}

	for _, page := range pages {
		// Search in filename
		found, score, preview := searchInLine(page.Filename, queries)
		if found {
			result.Filenames = append(result.Filenames, searchResult{
				Filename: page.Filename,
				ModTime:  page.ModTime,
				Score:    score,
				Preview:  preview,
			})
		}

		// Search in content
		markdownContent, _ := s.PageRawContent(page.Filename)
		content := stripmd.Strip(markdownContent)
		lines := strings.Split(content, "\n")

		r := searchResult{
			Filename: page.Filename,
			ModTime:  page.ModTime,
			Score:    0,
			Preview:  "",
		}

		for _, line := range lines {
			found, score, preview := searchInLine(line, queries)
			if found {
				r.Score += score
				r.Preview += "\n" + preview
			}
		}

		if r.Score > 0 {
			result.PageContent = append(result.PageContent, r)
		}
	}

	sort.Slice(result.PageContent, func(i, j int) bool {
		return result.PageContent[i].Score > result.PageContent[j].Score
	})

	sort.Slice(result.Filenames, func(i, j int) bool {
		return result.Filenames[i].Score > result.Filenames[j].Score
	})

	return &result, nil
}

func searchInLine(line string, queries []string) (bool, int, string) {
	score := 0
	preview := line
	words := strings.Fields(line)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	for _, query := range queries {
		for _, word := range words {
			tq, _, _ := transform.String(t, strings.ToLower(query))
			tw, _, _ := transform.String(t, strings.ToLower(word))

			index := strings.Index(tw, tq)
			if index != -1 {
				score += 1
				preview = highlight(preview, word, index, query)
			}
		}
	}

	if score > 0 {
		return true, score, preview
	}
	return false, score, ""
}

func extractQueries(query string) []string {
	out := make([]string, 0)
	for _, q := range strings.Fields(query) {
		out = append(out, strings.TrimSpace(q))
	}
	return out
}

// highlight 'word' in 'src' starting from 'index' to length of the 'query'
func highlight(src string, word string, index int, query string) string {
	r := word[:index] + "<b class=\"highlight\">" + word[index:index+len(query)] + "</b>" + word[index+len(query):]
	return strings.Replace(src, word, r, 1)
}
