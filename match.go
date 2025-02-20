package main

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"regexp"
	"strings"
)

type matchUsers struct {
	Adds    sets.String
	Removes sets.String
}

func (mu *matchUsers) isMatched() bool {
	if mu == nil {
		return false
	}
	return len(mu.Adds) != 0 || len(mu.Removes) != 0
}

func matchAssign(comment, commenter string) *matchUsers {
	assignRe := regexp.MustCompile(`(?mi)^/(un)?assign(( @?[-\w]+?)*)\s*$`)

	matches := assignRe.FindAllStringSubmatch(comment, -1)
	if len(matches) == 0 {
		return nil
	}
	isAdd := func(re string) bool {
		return re != "un"
	}

	return extractMatchUsers(commenter, matches, isAdd)
}

func matchCollaborator(comment, commenter string) *matchUsers {
	collaboratorRe := regexp.MustCompile(`(?mi)^/(add|rm)-collaborator(( @?[-\w]+?)*)\s*$`)

	matches := collaboratorRe.FindAllStringSubmatch(comment, -1)
	if matches == nil {
		return nil
	}

	isAdd := func(re string) bool {
		return re != "rm"
	}

	return extractMatchUsers(commenter, matches, isAdd)
}

func extractMatchUsers(commenter string, matches [][]string, isAdd func(string) bool) *matchUsers {
	toAdd := sets.NewString()
	toRemove := sets.NewString()
	save := func(login string, isAdd bool) {
		if isAdd {
			toAdd.Insert(login)
		} else {
			toRemove.Insert(login)
		}
	}

	for _, re := range matches {
		add := isAdd(re[1])
		if re[2] == "" {
			save(commenter, add)
		} else {
			for _, login := range parseLogins(re[2]) {
				save(login, add)
			}
		}
	}

	return &matchUsers{
		Adds:    toAdd,
		Removes: toRemove,
	}
}

func parseLogins(text string) []string {
	var parts []string
	for _, p := range strings.Split(text, " ") {
		t := strings.Trim(p, "@ ")
		if t == "" {
			continue
		}
		parts = append(parts, t)
	}
	return parts
}
