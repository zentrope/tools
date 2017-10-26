//
// Copyright (C) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type field struct {
	Name      string  `json:"name"`
	Delimiter string  `json:"delimiter"`
	Scrub     string  `json:"scrub"`
	Fields    []field `json:"fields"`
}

type grammer struct {
	Name               string  `json:"name"`
	OS                 string  `json:"os"`
	Delimiter          string  `json:"delimiter"`
	CondenseWhitespace bool    `json:"condenseWhitespace"`
	CheckContinuations bool    `json:"checkContinuations"`
	NumberOfFields     int     `json:"numberOfFields"`
	Fields             []field `json:"fields"`
}

func newGrammer(filename string) (*grammer, error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var grammer grammer
	err = json.Unmarshal(content, &grammer)
	if err != nil {
		return nil, err
	}

	return &grammer, nil
}

type logFact struct {
	Month       string
	Day         string
	Hour        string
	Minute      string
	Second      string
	ProcessID   string
	ProcessName string
	Message     string
}

func newLogFact(env fieldMap) logFact {
	return logFact{
		Month:       env["month"],
		Day:         env["day"],
		Hour:        env["hour"],
		Minute:      env["minute"],
		Second:      env["second"],
		ProcessID:   env["process-id"],
		ProcessName: env["process-name"],
		Message:     env["message"],
	}
}

// All logs assume a continuation is a line starting with whitespace
var continuationRE = regexp.MustCompile(`^\s`)

// All logs might want to condense whitespace.
var condenseRE = regexp.MustCompile(`\s`)

type fieldMap map[string]string

func (r field) parse(env fieldMap, line string) {
	if len(r.Fields) == 0 {
		env[r.Name] = line
		return
	}

	line2 := line
	tokensRE := regexp.MustCompile(r.Delimiter)
	if r.Scrub != "" {
		scrubRE := regexp.MustCompile(r.Scrub)
		line2 = scrubRE.ReplaceAllLiteralString(line2, "")
	}
	tokens := tokensRE.Split(line2, -1)
	for i, field := range r.Fields {
		field.parse(env, tokens[i])
	}
}

func (r *grammer) parse(line string) fieldMap {
	tokensRE := regexp.MustCompile(r.Delimiter)
	tokens := tokensRE.Split(line, r.NumberOfFields)

	data := make(fieldMap, 0)
	for i, field := range r.Fields {
		field.parse(data, tokens[i])
	}

	return data
}

func (r *grammer) isContinuation(line string) bool {
	if r.CheckContinuations {
		return continuationRE.MatchString(line)
	}
	return false
}

func (r *grammer) doOutput(line string) {
	if line == "" {
		return
	}

	line2 := line
	if r.CondenseWhitespace {
		line2 = condenseRE.ReplaceAllString(line2, " ")
		line2 = strings.TrimSpace(line2)
	}

	fmt.Printf("%#v\n", newLogFact(r.parse(line2)))
}

//-----------------------------------------------------------------------------
// MAIN
//-----------------------------------------------------------------------------

func main() {
	fmt.Println("Log Ripper")

	rules, err := newGrammer("grammar/syslog.json")
	if err != nil {
		log.Fatal(err)
	}

	input := bufio.NewReader(os.Stdin)

	var bigLine string

	for {
		line, err := input.ReadString('\n')

		// skip these for now
		if strings.Index(line, "--- last message repeated") != -1 {
			continue
		}

		if err != nil {
			rules.doOutput(bigLine)
			break
		}

		if !rules.isContinuation(line) {
			rules.doOutput(bigLine)
			bigLine = line
		} else {
			bigLine = bigLine + " " + strings.TrimSpace(line)
		}
	}
}
