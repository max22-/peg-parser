package main

type ParseResult interface {
	isParseResult()
}

type Parser func([]byte, int) (ParseResult, int)

type Char byte
type Int int
type List []ParseResult

func (c Char) isParseResult()   {}
func (i Int) isParseResult()    {}
func (prs List) isParseResult() {}

func char(c byte) Parser {
	return func(source []byte, loc int) (ParseResult, int) {
		if loc >= len(source) {
			return nil, loc
		}
		if c == source[loc] {
			return Char(c), loc + 1
		}
		return nil, loc
	}
}

func digit() Parser {
	return func(source []byte, loc int) (ParseResult, int) {
		if loc >= len(source) {
			return nil, loc
		}
		c := source[loc]
		if c >= '0' && c <= '9' {
			return Int(c - '0'), loc + 1
		}
		return nil, loc
	}
}

func seq(ps []Parser) Parser {
	return func(source []byte, loc int) (ParseResult, int) {
		if loc >= len(source) {
			return nil, loc
		}
		loc2 := loc
		var res []ParseResult
		for _, p := range ps {
			var r ParseResult
			r, loc2 = p(source, loc2)
			if r == nil {
				return nil, loc
			}
			res = append(res, r)
		}
		return List(res), loc2
	}
}

func choice(ps []Parser) Parser {
	return func(source []byte, loc int) (ParseResult, int) {
		if loc >= len(source) {
			return nil, loc
		}
		for _, p := range ps {
			res, loc2 := p(source, loc)
			if res != nil {
				return res, loc2
			}
		}
		return nil, loc
	}
}

func many(p Parser) Parser {
	return func(source []byte, loc int) (ParseResult, int) {
		var res []ParseResult
		pr, loc := p(source, loc)
		for pr != nil {
			res = append(res, pr)
			pr, loc = p(source, loc)
		}
		return List(res), loc
	}
}
