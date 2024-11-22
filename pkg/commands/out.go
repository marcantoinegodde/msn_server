package commands

func HandleOUT(c chan string) {
	res := "OUT\r\n"
	c <- res
}
