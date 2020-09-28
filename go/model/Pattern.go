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

	"github.com/antchfx/xmlquery"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-samples/cloud-suitability-analyzer/go/util"
	"gopkg.in/yaml.v3"
)

const ANALYZED_FILE_PATTERN = "Analyzed File"
const FILE_SLOC_PATTERN = "Lines of Code"

type Pattern struct {
	ID            uint           `gorm:"primary_key" json:"-" yaml:"-"`
	CreatedAt     time.Time      `json:"-" yaml:"-"`
	UpdatedAt     time.Time      `json:"-" yaml:"-"`
	RuleID        uint           `sql:"type:bigint REFERENCES rules(id) ON DELETE CASCADE" json:"-"  yaml:"-"`
	Rule          Rule           `gorm:"foreignkey:RuleID" json:"-"  yaml:"-"`
	Type          string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"` //Empty if not overriding base rule type
	Pattern       string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"` //Empty if not overriding base rule regex
	Value         string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"`
	Advice        string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"`
	Effort        int            `gorm:"type:bigint" json:"effort,omitempty" yaml:"effort,omitempty"`
	Readiness     int            `gorm:"type:bigint" json:"readiness,omitempty" yaml:"readiness,omitempty"`
	Criticality   string         `json:"criticality,omitempty" yaml:"criticality,omitempty"`
	Tag           string         `gorm:"index;type:text" json:"tag,omitempty" yaml:"tag,omitempty"`
	Recipe        string         `gorm:"index;type:text" json:"recipe,omitempty" yaml:"recipe,omitempty"`
	Category      string         `json:"category,omitempty" yaml:"category,omitempty"`
	compiledRegex *regexp.Regexp `gorm:"-" json:"-" yaml:"-"`
	Command       string         `gorm:"type:text" json:",omitempty" yaml:",omitempty"`
	sync.Mutex    `gorm:"-" json:"-" yaml:"-"`
}

func (p *Pattern) IsValid(rule *Rule) error {

	if len(p.Value) < 1 && len(p.Command) < 1 {
		return fmt.Errorf("A valid rule pattern requires a value or command!")
	}

	if p.Type == REGEX_MATCH_TYPE || rule.Type == REGEX_MATCH_TYPE {
		if p.Pattern != "" {
			patt := strings.Replace(p.Pattern, "%s", "submarker", 1)
			_, err := regexp.Compile(patt)

			if err != nil {
				return fmt.Errorf("regex pattern [%s] for pattern value [%s] is bad. Details: %s", p.Pattern, p.Value, err.Error())
			}
		}
	}
	return nil
}

func (p *Pattern) ToString() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("\tType: %s", p.Type))
	b.WriteString(fmt.Sprintf("\tPattern: %s", p.Pattern))
	b.WriteString(fmt.Sprintf("\tValue: %s", p.Value))
	b.WriteString(fmt.Sprintf("\tAdvice: %s", p.Advice))
	b.WriteString(fmt.Sprintf("\tEffort: %d", p.Effort))
	b.WriteString(fmt.Sprintf("\tReadiness: %d", p.Readiness))
	b.WriteString(fmt.Sprintf("\tTag: %s", p.Tag))
	b.WriteString(fmt.Sprintf("\tCategory: %s", p.Category))
	b.WriteString(fmt.Sprintf("\tRecipe: %s", p.Recipe))

	return b.String()
}

func (p *Pattern) GetEscapedValue() string {
	return util.EscapeSpecials(p.Value)
}

func (p *Pattern) GetEscapedPattern() string {
	return util.EscapeSpecials(p.Pattern)
}

func (p *Pattern) compile(rule *Rule) {

	//Rule or Pattern Governed match
	if p.Type == "" {
		p.Type = rule.Type
	}

	if p.Pattern == "" {
		p.Pattern = rule.DefaultPattern
	}

	if strings.Contains(p.Pattern, "%s") {
		p.Pattern = strings.Replace(p.Pattern, "%s", p.Value, 1)
	} else {
		p.Pattern = p.Value
	}

	if p.Type == REGEX_MATCH_TYPE {
		p.compiledRegex = regexp.MustCompile(p.Pattern)
	}
}

func (p *Pattern) MatchYaml(node *yaml.Node) (bool, string) {
	if node != nil {
		switch p.Type {
		case YAMLPATH_MATCH_TYPE:
			if path, err := yamlpath.NewPath(p.Pattern); err == nil {
				if results, err := path.Find(node); err == nil {
					return (len(results) > 0), ""
				}
			}
		}
	}

	return false, ""
}

func (p *Pattern) MatchXml(node *xmlquery.Node) (bool, string) {
	switch p.Type {
	case XPATH_MATCH_TYPE:
		return (xmlquery.FindOne(node, p.Pattern) != nil), ""
	}

	return false, ""
}

func (p *Pattern) Match(target string) (bool, string) {

	switch p.Type {

	case REGEX_MATCH_TYPE:
		if p.compiledRegex.MatchString(target) {
			return true, ""
		}
	case STARTS_WITH_CI_MATCH_TYPE:
		if strings.HasPrefix(strings.ToLower(target), strings.ToLower(p.Pattern)) {
			return true, ""
		}
	case STARTS_WITH_MATCH_TYPE:
		if strings.HasPrefix(target, p.Pattern) {
			return true, ""
		}
	case ENDS_WITH_CI_MATCH_TYPE:
		if strings.HasSuffix(strings.ToLower(target), strings.ToLower(p.Pattern)) {
			return true, ""
		}
	case ENDS_WITH_MATCH_TYPE:
		if strings.HasSuffix(target, p.Pattern) {
			return true, ""
		}
	case CONTAINS_CI_MATCH_TYPE:
		if strings.Contains(strings.ToLower(target), strings.ToLower(p.Pattern)) {
			return true, ""
		}
	case CONTAINS_MATCH_TYPE:
		if strings.Contains(target, p.Pattern) {
			return true, ""
		}
	case SIMPLE_TEXT_CI_MATCH_TYPE:
		if strings.ToLower(target) == strings.ToLower(p.Pattern) {
			return true, ""
		}
	default:
		if target == p.Pattern {
			return true, ""
		}
	}

	return false, ""
}
