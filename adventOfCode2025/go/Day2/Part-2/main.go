package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func countInvalidIDs() int {
	file, err := os.Open("./adventOfCode2025/go/Day2/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		pairs := strings.Split(line, ",")

		for _, pair := range pairs {
			ids := strings.Split(pair, "-")
			if len(ids) != 2 {
				continue
			}

			start, err := strconv.Atoi(ids[0])
			if err != nil {
				log.Fatalf("failed to convert start ID to integer: %s", err)
			}

			end, err := strconv.Atoi(ids[1])
			if err != nil {
				log.Fatalf("failed to convert end ID to integer: %s", err)
			}

			for id := start; id <= end; id++ {
				idStr := strconv.Itoa(id)
				idConcat := idStr + idStr
				idRotated := idConcat[1 : len(idConcat)-1]
				if strings.Contains(idRotated, idStr) {
					sum += id
				}
			}
		}
	}

	return sum
}

func main() {
	result := countInvalidIDs()
	fmt.Println("Result:", result)
}

/*
Let me make a note for future reference.

If the string S has repeated block, it could be described in terms of pattern.
S = SpSp (For example, S has two repeatable block at most)
If we repeat the string, then SS=SpSpSpSp.
Destroying first and the last pattern by removing each character, we generate a new S2=SxSpSpSy.

If the string has repeatable pattern inside, S2 should have valid S in its string.

The maximum length of a "repeated" substring that you could get from a string would be half it's length
For example, s = "abcdabcd", "abcd" of len = 4, is the repeated substring.
You cannot have a substring >(len(s)/2), that can be repeated.

So, when ss = s + s , we will have atleast 4 parts of "repeated substring" in ss.
(s+s)[1:-1], With this we are removing 1st char and last char => Out of 4 parts of repeated substring, 2 part will be gone (they will no longer have the same substring).
ss.find(s) != -1, But still we have 2 parts out of which we can make s. And that's how ss should have s, if s has repeated substring.
*/
