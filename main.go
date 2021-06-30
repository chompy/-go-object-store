package main

//var client gokv.Store

func main() {
	if err := listen(8081); err != nil {
		panic(err)
	}
}
