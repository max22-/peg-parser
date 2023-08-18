package main

import "fmt"

func main() {
	stringify := func(l ParseResult[[]byte]) ParseResult[string] {
		var res ParseResult[string]
		res.success = false
		if l.success {
			res.val = string(l.val)
			res.success = true
		}
		return res
	}
	abc := apply(stringify, seq([]Parser[byte]{char('a'), char('b'), char('c')}))
	pr, loc := many(abc)([]byte("abcabc123abc"), 0)
	fmt.Println(pr, loc)
}
