package specter

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFromStdin() {
	// Receive the input via stdin
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stdin: ", err)
	}
}
