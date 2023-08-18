package main

import "fmt"

func main() {
	abc := seq([]Parser{char('a'), char('b'), char('c')})
	pr, loc := many(abc)([]byte("abcabc123abc"), 0)
	fmt.Println(pr, loc)
}
