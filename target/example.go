package target

import "fmt"

func main() {
	fmt.Println("Start")

	a := 10
	b := 20

	if a < b {
		fmt.Println("a is less than b")
	} else {
		fmt.Println("a is not less than b")
	}

	for i := 0; i < 3; i++ {
		fmt.Println("Loop iteration", i)
	}

	fmt.Println("End")
}
