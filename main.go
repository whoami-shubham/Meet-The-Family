package main

func main() {
	build := "build.txt"
	test := "test.txt"
	commands(build, true) // build family
	commands(test, false) // sample test
}
