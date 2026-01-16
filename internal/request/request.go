package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	header, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	splittedHeader := strings.Split(string(header), "\r\n")

	requestLine, _, err := parseRequestLine(splittedHeader[0])
	if err != nil {
		return nil, err
	}

	request := Request{*requestLine}

	return &request, nil
}

func parseRequestLine(requestLine string) (*RequestLine, int, error) {
	headerParts := strings.Split(requestLine, " ")
	if len(headerParts) != 3 {
		requestLineLenErr := errors.New("Bad request line")
		return nil, 0, requestLineLenErr
	}

	methodRegex, err := regexp.Compile(`^[A-Z]+$`)
	if err != nil {
		return nil, 0, err
	}

	method := headerParts[0]
	if !methodRegex.Match([]byte(method)) {
		methodErr := errors.New("Bad HTTP method")
		return nil, 0, methodErr
	}

	targetRegex, err := regexp.Compile(`^/([\w\-]+/?)*$`)
	if err != nil {
		return nil, 0, err
	}

	target := headerParts[1]
	if !targetRegex.Match([]byte(target)) {
		targetErr := errors.New("Bad HTTP target")
		return nil, 0, targetErr
	}

	version := headerParts[2]
	versionParts := strings.Split(version, "/")

	if len(versionParts) != 2 || versionParts[0] != "HTTP" {
		versionErr := errors.New("Bad HTTP version")
		return nil, 0, versionErr
	}

	versionNumber := versionParts[1]
	if versionNumber != "1.1" {
		versionNumberErr := errors.New("Unsupported HTTP version")
		return nil, 0, versionNumberErr
	}

	parsedReqeuestLine := RequestLine{versionNumber, target, method}

	return &parsedReqeuestLine, 0, nil
}
