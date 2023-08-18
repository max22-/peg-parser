package main

import "fmt"

func main() {
	abc := apply(func(bs []byte) string { return string(bs) }, seq([]Parser[byte]{char('a'), char('b'), char('c')}))
	pr, loc := many(abc)([]byte("abcabc123abc"), 0)
	fmt.Println(pr, loc)
}
