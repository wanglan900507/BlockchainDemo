package main


import "fmt"


func main() {
	var i, j int = 1, 2
	k := 3

	fmt.Println(i, j, k)

	fmt.Println(split(20))

	fmt.Println(transferAToB(100, 200, 10))
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func transferAToB(a, b, transferValue int) (aAfter, bAfter int){
	aAfter = a - transferValue
	bAfter = b + transferValue

	return
}



