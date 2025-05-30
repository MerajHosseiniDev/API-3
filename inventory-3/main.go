package main

func main() {
	app := App{}
	app.Initialise()
	app.Run("localhost:6666")
}