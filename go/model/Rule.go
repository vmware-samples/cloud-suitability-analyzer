/*******************************************************************************
 * Copyright (c) 2018 - Present VMWare, Inc. All Rights Reserved.
 * SPDX-License-Identifier: BSD-2
 ******************************************************************************/

package model

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/vmware-samples/cloud-suitability-analyzer/go/util"
)

type Rule struct {
	ID              uint           `gorm:"primary_key" json:"-" yaml:"-"`
	CreatedAt       time.Time      `json:"-" yaml:"-"`
	UpdatedAt       time.Time      `json:"-" yaml:"-"`
	Name            string         `gorm:"type:text;unique_index;not null"`
	FileType        string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"` //Extension if empty or * then rule applies to all files. Else, rule only applies to files with this extension (sans '.')
	Target          string         `gorm:"type:text"`                                     //File, Line
	Type            string         `gorm:"type:text"`                                     //Regex, SimpleText, StartsWith, Contains, EndsWith, SimpleTextCaseInsensitive
	DefaultPattern  string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"` //Pattern with a placeholder that follows standard GO fmt.Sprintf rules. I.E. "[ .]%s[ (]" uses %s for string Pattern Value substitution before compilation!
	Advice          string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"`
	Effort          int            `gorm:"type:bigint; column:effort" json:",omitempty" yaml:",omitempty"`
	Impact          string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"`
	Readiness       int            `gorm:"type:bigint; column:readiness" json:",omitempty" yaml:",omitempty"`
	Category        string         `json:",omitempty" yaml:",omitempty"`
	Criticality     string         `json:",omitempty" yaml:",omitempty"`
	Tags            []Tag          `json:",omitempty" yaml:",omitempty"`
	Recipes         []Recipe       `gorm:"foreignkey:RuleID" json:",omitempty" yaml:",omitempty"`
	Patterns        []Pattern      `gorm:"foreignkey:RuleID"`
	regex           *regexp.Regexp `gorm:"-" json:"-" yaml:"-"`
	Metric          *RuleMetric    `gorm:"-" json:"-" yaml:"-"`
	overrideApplies bool           `gorm:"-" json:"-" yaml:"-"`
	sync.Mutex      `gorm:"-" json:"-" yaml:"-"`
}

func (r *Rule) Applies(fileExt string) bool {
	if *util.Verbose {
		fmt.Printf("Checking to see if rule [%s] applies to extension [%s] with pattern [%s]\n", r.Name, fileExt, r.FileType)
	}

	if r.overrideApplies {
		return true
	}

	return r.regex.MatchString(fileExt)
}

func (r *Rule) IsValid() (isValid bool, err error) {

	//Name is required
	if r.Name == "" {
		return false, fmt.Errorf("Rule Name is required!")
	}

	//FileType
	if r.FileType != "" && r.FileType != "*" {
		//Check FileType is valid regex!
		_, err = regexp.Compile(r.FileType)
		if err != nil {
			return false, fmt.Errorf("FileType must be a empty (nil, not specified), or \"*\", or valid regex! FileType: %s Fails with error: %v", r.FileType, err)
		}
	}

	//Target
	if r.Target == "" {
		return false, fmt.Errorf("Rule Target is required!")
	}

	if r.Target != FILE_TARGET && r.Target != LINE_TARGET && r.Target != CONTENTS_TARGET {
		return false, fmt.Errorf("Rule Target must be either (%s|%s|%s) but was: %s",
			FILE_TARGET, LINE_TARGET, CONTENTS_TARGET, r.Target)
	}

	//Type
	isValid, err = r.typeIsValid()
	if !isValid {
		return isValid, err
	}

	isValid, err = r.impactIsValid()
	if !isValid {
		return isValid, err
	}

	//Patterns
	if len(r.Patterns) < 1 {
		return false, fmt.Errorf("Rule must have at least one (1) pattern")
	} else {
		for _, pattern := range r.Patterns {
			err = pattern.IsValid(r)
			if err != nil {
				return false, fmt.Errorf("Rule Pattern [%s] Invalid! Details: %s", pattern.ToString(), err.Error())
			}
		}
	}

	isValid, err = r.defaultPatternIsValid()
	if !isValid {
		return isValid, err
	}

	return true, nil
}

func (r *Rule) AddPattern(pattern Pattern) {
	//Is that really all there is to it?

	r.Patterns = append(r.Patterns, pattern)
	if *util.Verbose {
		fmt.Printf("Adding Pattern [%s] to rule [%s] new pattern cnt: %d\n", pattern.Value, r.Name, len(r.Patterns))
	}
}

func (r *Rule) AddRecipe(recipe Recipe) {
	//Is that really all there is to it?

	r.Recipes = append(r.Recipes, recipe)
	if *util.Verbose {
		fmt.Printf("Adding Recipe [%s] to rule [%s] new Recipe cnt: %d\n", recipe.URI, r.Name, len(r.Recipes))
	}
}

func (r *Rule) AddTag(tag string) {
	r.Tags = append(r.Tags, Tag{Value: tag})
}

func (r *Rule) HasTag(tag string) bool {
	if tag != "" && len(r.Tags) > 0 {
		for i := range r.Tags {
			//Case insensitive matching!
			if strings.ToLower(tag) == strings.ToLower(r.Tags[i].Value) {
				return true
			}
		}
	}
	return false
}

func (r *Rule) CompilePatterns() {
	r.Lock()
	defer r.Unlock()

	if r.FileType == "" || r.FileType == "*" {
		r.overrideApplies = true
	}

	//Compile the rule file type regex!
	r.regex = regexp.MustCompile(r.FileType)

	for i, _ := range r.Patterns {
		r.Patterns[i].compile(r)
	}
}

func (r *Rule) GetEscapedPattern() string {
	return util.EscapeSpecials(r.DefaultPattern)
}

func (r *Rule) UpdateRule(newRule Rule) (deletedPatterns []Pattern, deletedRecipes []Recipe, deletedTags []Tag) {

	if newRule.Advice != "" && newRule.Advice != r.Advice {
		r.Advice = newRule.Advice
	}

	if newRule.Effort != -1 && newRule.Effort != r.Effort {
		r.Effort = newRule.Effort
	}

	if newRule.Impact != "" && newRule.Impact != r.Impact {
		r.Impact = newRule.Impact
	}

	if newRule.Target != "" && newRule.Target != r.Target {
		r.Target = newRule.Target
	}

	if newRule.Category != "" && newRule.Category != r.Category {
		r.Category = newRule.Category
	}

	if newRule.FileType != "" && newRule.FileType != r.FileType {
		r.FileType = newRule.FileType
	}

	if newRule.Type != "" && newRule.Type != r.Type {
		r.Type = newRule.Type
	}

	if newRule.DefaultPattern != "" && newRule.DefaultPattern != r.DefaultPattern {
		r.DefaultPattern = newRule.DefaultPattern
	}

	if newRule.Criticality != "" && newRule.Criticality != r.Criticality {
		r.Criticality = newRule.Criticality
	}

	deletedPatterns = r.updatePatterns(newRule)
	deletedRecipes = r.updateRecipes(newRule)
	deletedTags = r.updateTags(newRule)

	return
}

func (r *Rule) updatePatterns(newRule Rule) (deleted []Pattern) {

	var newPatterns int
	var updatedPatterns int

	if len(newRule.Patterns) > 0 {

		for _, pattern := range newRule.Patterns {

			var patternMatched = false
			var patternUpdated = false

			for i, _ := range r.Patterns {

				if pattern.Value == r.Patterns[i].Value {
					//update pattern
					patternMatched = true

					if pattern.Advice != "" && pattern.Advice != r.Patterns[i].Advice {
						r.Patterns[i].Advice = pattern.Advice
						patternUpdated = true
					}

					if pattern.Effort != -1 && pattern.Effort != r.Patterns[i].Effort {
						r.Patterns[i].Effort = pattern.Effort
						patternUpdated = true
					}

					if pattern.Readiness != -1 && pattern.Readiness != r.Patterns[i].Readiness {
						r.Patterns[i].Readiness = pattern.Readiness
						patternUpdated = true
					}

					if pattern.Type != "" && pattern.Type != r.Patterns[i].Type {
						r.Patterns[i].Type = pattern.Type
						patternUpdated = true
					}

					if pattern.Pattern != "" && pattern.Pattern != r.Patterns[i].Pattern {
						r.Patterns[i].Pattern = pattern.Pattern
						patternUpdated = true
					}
					if patternUpdated {
						updatedPatterns++
					}
					break
				}
			}

			if !patternMatched {
				r.Patterns = append(r.Patterns, pattern)
				newPatterns++
			}

		}

		for i := len(r.Patterns) - 1; i >= 0; i-- {
			var found = false
			existingPattern := r.Patterns[i]

			for _, newPattern := range newRule.Patterns {
				if existingPattern.Value == newPattern.Value {
					found = true
					break
				}
			}

			if !found {
				r.Patterns = append(r.Patterns[:i], r.Patterns[i+1:]...)
				deleted = append(deleted, existingPattern)
			}
		}
		if *util.Verbose {
			fmt.Printf(" Patterns[updated:%d added:%d deleted:%d] ", updatedPatterns, newPatterns, len(deleted))
		}
	}

	return deleted
}

func (r *Rule) updateRecipes(newRule Rule) (deleted []Recipe) {

	var newRecipes int

	if len(newRule.Recipes) > 0 {
		for _, recipe := range newRule.Recipes {
			var recipeMatched = false
			for i, _ := range r.Recipes {

				if recipe.URI == r.Recipes[i].URI {
					//update recipe
					recipeMatched = true
				}
			}

			if !recipeMatched {
				r.Recipes = append(r.Recipes, recipe)
				newRecipes++
			}
		}

		for i := len(r.Recipes) - 1; i >= 0; i-- {

			var found = false

			recipe := r.Recipes[i]

			for _, matchedRecipe := range newRule.Recipes {
				if recipe.URI == matchedRecipe.URI {
					found = true
					break
				}
			}

			if !found {
				r.Recipes = append(r.Recipes[:i], r.Recipes[i+1:]...)
				deleted = append(deleted, recipe)
			}
		}

		if *util.Verbose {
			fmt.Printf(" Recipes[added:%d deleted:%d] ", newRecipes, len(deleted))
		}
	}

	return deleted
}

func (r *Rule) updateTags(newRule Rule) (deleted []Tag) {

	var newTags int

	if len(newRule.Tags) > 0 {

		for _, tag := range newRule.Tags {

			var tagMatched = false

			for i, _ := range r.Tags {

				if tag.Value == r.Tags[i].Value {
					tagMatched = true
				}
			}

			if !tagMatched {
				r.Tags = append(r.Tags, tag)
				newTags++
			}
		}

		for i := len(r.Tags) - 1; i >= 0; i-- {

			var found = false

			tag := r.Tags[i]

			for _, matchedTag := range newRule.Tags {
				if tag.Value == matchedTag.Value {
					found = true
					break
				}
			}

			if !found {
				r.Tags = append(r.Tags[:i], r.Tags[i+1:]...)
				deleted = append(deleted, tag)
			}
		}

		if *util.Verbose {
			fmt.Printf(" Tags[added:%d deleted:%d] ", newTags, len(deleted))
		}
	}

	return deleted
}

func GetRulesDir(isImport bool) string {

	dir := *util.RulesDir

	//User overrode default!
	if dir != util.DEFAULT_RULES_DIR || isImport {
		return dir
	}

	dir = *util.OutputDir
	if !strings.HasSuffix(dir, util.PathSeparator) {
		dir = dir + util.PathSeparator
	}

	return dir + "rules"

}

func (r *Rule) typeIsValid() (isValid bool, err error) {
	//Type
	if r.Type == "" {
		return false, fmt.Errorf("Rule Type is required!")
	}

	if r.Type != REGEX_MATCH_TYPE &&
		r.Type != XPATH_MATCH_TYPE &&
		r.Type != YAMLPATH_MATCH_TYPE &&
		r.Type != PLUGIN_MATCH_TYPE &&
		r.Type != SIMPLE_TEXT_MATCH_TYPE &&
		r.Type != SIMPLE_TEXT_CI_MATCH_TYPE &&
		r.Type != STARTS_WITH_MATCH_TYPE &&
		r.Type != STARTS_WITH_CI_MATCH_TYPE &&
		r.Type != ENDS_WITH_MATCH_TYPE &&
		r.Type != ENDS_WITH_CI_MATCH_TYPE &&
		r.Type != CONTAINS_MATCH_TYPE &&
		r.Type != CONTAINS_CI_MATCH_TYPE {
		return false, fmt.Errorf("Rule Type must be (%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s) but was: %s",
			REGEX_MATCH_TYPE,
			XPATH_MATCH_TYPE,
			YAMLPATH_MATCH_TYPE,
			PLUGIN_MATCH_TYPE,
			SIMPLE_TEXT_CI_MATCH_TYPE,
			SIMPLE_TEXT_MATCH_TYPE,
			STARTS_WITH_CI_MATCH_TYPE,
			STARTS_WITH_MATCH_TYPE,
			ENDS_WITH_CI_MATCH_TYPE,
			ENDS_WITH_MATCH_TYPE,
			CONTAINS_CI_MATCH_TYPE,
			CONTAINS_MATCH_TYPE, r.Target)
	}

	return true, nil

}

func (r *Rule) impactIsValid() (isValid bool, err error) {
	//Impact is not required
	if r.Impact == "" {
		return true, nil
	}

	if r.Impact != EVERY_IMPACT &&
		r.Impact != FILE_IMPACT &&
		r.Impact != APP_IMPACT {
		return false, fmt.Errorf("Rule Impact must be (%s|%s|%s) but was: %s",
			EVERY_IMPACT,
			FILE_IMPACT,
			APP_IMPACT,
			r.Impact)
	}

	return true, nil

}

func (r *Rule) defaultPatternIsValid() (isValid bool, err error) {

	//Contains substitution for patterns...
	if r.DefaultPattern != "" && !strings.Contains(r.DefaultPattern, "%s") {
		return false, fmt.Errorf("If provided Rule DefaultPattern should include a substitution marker of '%%s'! Actual: %s", r.DefaultPattern)
	}

	if r.Type == REGEX_MATCH_TYPE && r.DefaultPattern != "" {

		//Substitute a value!
		pattern := strings.Replace(r.DefaultPattern, "%s", "submarker", 1)

		_, err = regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Errorf("DefaultPattern for '%s' rules must be a valid perl compatible regular expression! Pattern: %s Fails with error: %v", REGEX_MATCH_TYPE, r.DefaultPattern, err)
		}

	}

	if r.DefaultPattern == "" {
		//If I don't have a default pattern I must contain at least one pattern with a pattern!
		if len(r.Patterns) < 1 {
			return false, fmt.Errorf("Rule has no default pattern and no patterns!")
		}
	}

	return true, nil
}
