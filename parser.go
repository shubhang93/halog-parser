package main

import (
	"errors"
	"regexp"
	"strings"
)

type IncomingCall struct {
	Backend string `json:"backend"`
	URL     string `json:"url"`
	Verb    string `json:"verb"`
}

func parseLine(line string, rgx *regexp.Regexp) (IncomingCall, error) {
	index := strings.Index(line, subStr)

	if index < 0 {
		return IncomingCall{}, errors.New("could not find substr:" + subStr)
	}

	rest := line[index+len(subStr)+1:]
	backendNodeName := ""
	for i := 0; i < len(rest); i++ {
		if string(rest[i]) == " " {
			break
		}
		backendNodeName += string(rest[i])
	}

	line = line[index:]

	httpReqLog := strings.TrimSpace(rgx.FindString(line))
	if httpReqLog == "" {
		return IncomingCall{}, errors.New("cannot parse httpReqLog for line: " + line)
	}
	requestComps := strings.Split(httpReqLog, " ")
	if len(requestComps) > 2 {
		requestComps = requestComps[:2]
	}
	backend := strings.Split(backendNodeName, "/")[0]

	return IncomingCall{
		URL:     requestComps[1],
		Backend: backend,
		Verb:    requestComps[0],
	}, nil
}
