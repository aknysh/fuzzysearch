// Fuzzy searching allows for flexibly matching a string with partial input,
// useful for filtering data very quickly based on lightweight user input.
package fuzzy

import "unicode/utf8"

// Match returns true if source matches target using a fuzzy-searching
// algorithm. Note that it doesn't implement Levenshtein distance (see
// RankMatch instead), but rather a simplified version where there's no
// approximation. The method will return true only if each character in the
// source can be found in the target and occurs after the preceding matches.
func Match(source, target string) bool {
	return RankMatch(source, target) >= 0
}

// Find will return a list of strings in targets that fuzzy matches source.
func Find(source string, targets []string) []string {
	var matches []string

	for _, target := range targets {
		if Match(source, target) {
			matches = append(matches, target)
		}
	}

	return matches
}

// RankMatch is similar to Match except it will measure the Levenshtein
// distance between the source and the target and return its result. If there
// was no match, it will return -1.
func RankMatch(source, target string) int {
	lenDiff := len(target) - len(source)

	if lenDiff < 0 {
		return -1
	}

	if lenDiff == 0 && source == target {
		return 0
	}

	runeDiff := 0

Outer:
	for _, r1 := range source {
		for i, r2 := range target {
			if r1 == r2 {
				target = target[i+utf8.RuneLen(r2):]
				continue Outer
			} else {
				runeDiff++
			}
		}
		return -1
	}

	// count up remaining char
	for _, _ = range target {
		runeDiff++
	}

	return runeDiff
}

// RankFind is similar to Find, except it will also rank all matches using
// Levenshtein distance.
func RankFind(source string, targets []string) ranks {
	var r ranks
	for _, target := range Find(source, targets) {
		r = append(r, Rank{
			Source:   source,
			Target:   target,
			Distance: LevenshteinDistance(source, target),
		})
	}
	return r
}

type Rank struct {
	// Source is used as the source for matching.
	Source string
	// Target is the word matched against.
	Target string
	// Distance is the Levenshtein distance between Source and Target.
	Distance int
}

type ranks []Rank

func (r ranks) Len() int {
	return len(r)
}

func (r ranks) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ranks) Less(i, j int) bool {
	return r[i].Distance < r[j].Distance
}
