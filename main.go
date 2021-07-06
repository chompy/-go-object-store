package main

//var client gokv.Store

func main() {

	c, err := loadConfig()
	if err != nil {
		panic(err)
	}

	if err := listen(c); err != nil {
		panic(err)
	}
}
