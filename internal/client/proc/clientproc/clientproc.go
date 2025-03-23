package clientproc

import "fmt"

func RunProcess() error {
	fmt.Println("ClientProc run")
	defer fmt.Println("ClientProc end")

	return nil
}
