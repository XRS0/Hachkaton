package main

import (
	"fmt"
)

func main() {
	fmt.Println(solution("abc", "bc"))
	fmt.Println(solution("abc", "d"))
}

func solution(str, end string) bool {
	return len(str) >= len(end) && str[len(str)-len(end):] == end
}
