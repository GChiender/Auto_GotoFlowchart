package target

import "fmt"

func main() {
	name := "Go"
	age := 15

	fmt.Println("Hello,", name)

	if age > 10 {
		fmt.Println("Age is greater than 10")
	} else {
		fmt.Println("Age is 10 or less")
	}

	for i := 0; i < 3; i++ {
		fmt.Println("Loop:", i)
	}
}
