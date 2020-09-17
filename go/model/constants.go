/*******************************************************************************
 * Copyright (c) 2018 - Present VMWare, Inc. All Rights Reserved.
 * SPDX-License-Identifier: BSD-2
 ******************************************************************************/

package model

import "math"

const DATA_FIELD_PREFIX string = "Data"
const TOTAL_FIELD string = "***TOTAL"

const FILE_TARGET string = "file"         //Only matches against filenames!
const LINE_TARGET string = "line"         //Matches at a line level within a file matching the file ext!
const CONTENTS_TARGET string = "contents" //Matches against the entire files contents at one time!
const EVERY_IMPACT string = "every"       //A hit on this rule has an impact to effort/rediness for every match/occurence
const FILE_IMPACT string = "file"         //A hit on this rule has an impact only once per file. I.e. 10 findings in a file would only have the impact of 1 finding...the other 9 would be zero'd out
const APP_IMPACT string = "app"           //A hit on this rule has an impact only once per run regardless of the number of findings.
const REGEX_MATCH_TYPE string = "regex"
const XPATH_MATCH_TYPE string = "xpath"
const JSONPATH_MATCH_TYPE string = "jsonpath"
const YAMLPATH_MATCH_TYPE string = "yamlpath"
const PLUGIN_MATCH_TYPE string = "plugin"
const SIMPLE_TEXT_MATCH_TYPE string = "simple-text"
const SIMPLE_TEXT_CI_MATCH_TYPE string = "simple-text-ci"
const STARTS_WITH_MATCH_TYPE string = "starts-with"
const STARTS_WITH_CI_MATCH_TYPE string = "starts-with-ci"
const ENDS_WITH_MATCH_TYPE string = "ends-with"
const ENDS_WITH_CI_MATCH_TYPE string = "ends-with-ci"
const CONTAINS_MATCH_TYPE string = "contains"
const CONTAINS_CI_MATCH_TYPE string = "contains-ci"
const DEFAULT_CRITICALITY string = "medium"
const DEFAULT_LINE_REGEX_PATTERN string = "[ .]%s[ (]"
const THIRD_PARTY_TAG string = "third-party"
const IMPORTS_TAG string = "imports"
const ANNOTATION_TAG string = "annotations"
const API_TAG string = "api"
const FILE_TAG string = "file"
const MAX_RULE_SCORE int = 5
const CRIT_SCORE_THRESHOLD int = 9 //Findings with scores greater than this are considered critical
const MIN_RULE_SCORE int = 1

const THIRD_PARTY_REPORT_ID int = 1
const THIRD_PARTY_HEADER string = "ThirdParty"

const API_SUMMARY_REPORT_ID int = 2
const API_SUMMARY_API_HEADER string = "API"
const API_SUMMARY_API_COUNT_HEADER string = "Count"

const API_DETAILED_REPORT_ID int = 3
const API_DETAILED_DOMAIN_HEADER string = "Application"
const API_DETAILED_API_HEADER string = "API"
const API_DETAILED_PATTERN_HEADER string = "Patter"
const API_DETAILED_FILE_HEADER string = "File"
const API_DETAILED_LINE_NO_HEADER string = "LineNo"
const API_DETAILED_SOURCE_HEADER string = "Source"
const API_DETAILED_SCORE_HEADER string = "Score"
const API_DETAILED_ADVICE_HEADER string = "Advice"

const ANNOTATIONS_REPORT_ID int = 4
const ANNOTATIONS_HEADER string = "Annotations"

const CLOC_REPORT_ID int = 5
const CLOC_LANG_HEADER string = "language"
const CLOC_FILES_HEADER string = "#files"
const CLOC_BLANK_LINES_HEADER string = "blank"
const CLOC_COMMENT_LINES_HEADER string = "comment"
const CLOC_CODE_LINES_HEADER string = "code"

const THIRD_PARTY_IMPORTS string = "third-party-reports"
const THIRD_PARTY_IMPORTS_DESC string = "Lists 3rd Party Imports Used"
const API_SUMMARY string = "api-summary"
const API_SUMMARY_DESC string = "Summary of API(s) used with counts"
const API_DETAILED string = "api-detailed"
const API_DETAILED_DESC string = "Detailed listing of API(s) used with recommendations"
const ANNOTATIONS string = "annotations"
const ANNOTATIONS_DESC string = "Lists Annotations used"
const JAVA_2_ENGLISH string = "java-to-english"
const JAVA_2_ENGLISH_DESC string = "Plain English translation of java source for domain analysis"
const CLOC string = "cloc"
const CLOC_DESC string = "Lists source lines of code by type"

const GIT_FORENSICS_REPORT_ID int = 1
const GIT_FORENSICS string = "git-forensics"
const GIT_FORENSICS_DESC string = "Git forensics"
const GIT_ACTIVITY_DETAILED_REPORT_ID int = 2
const GIT_ACTIVITY_DETAILED string = "git-activity-detailed"
const GIT_ACTIVITY_DETAILED_DESC string = "Git detailed activity"
const GIT_ACTIVITY_SUMMARY_REPORT_ID int = 3
const GIT_ACTIVITY_SUMMARY string = "git-activity-summary"
const GIT_ACTIVITY_SUMMARY_DESC string = "Git activity summary"

const CSV_EXTENSION string = "csv"
const TXT_EXTENSION string = "txt"
const NO_ADVICE string = "No advice"

const All_criticality = "all"
const All_criticality_low_score int = math.MinInt64
const All_criticality_high_score int = math.MaxInt64
const Info_criticality = "info"
const Info_criticality_low_score int = math.MinInt64
const Info_criticality_high_score int = 0
const Low_criticality = "low"
const Low_criticality_low_score int = 1
const Low_criticality_high_score int = 3
const Medium_criticality = "medium"
const Medium_criticality_low_score int = 4
const Medium_criticality_high_score int = 6
const High_criticality = "high"
const High_criticality_low_score int = 7
const High_criticality_high_score int = math.MaxInt64
