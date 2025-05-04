package loggerf

import "fmt"

func Log(txt string, err error) {
	fmt.Println(txt+": %w", err)
}
