package main

func main() {
	for _, arrosage := range arrosage() {
		arrosage.arroser()
		arrosage.setDone().save()
	}
}
