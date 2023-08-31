package main

import "publisher/publisher"

func main() {
	pub := publisher.New()
	pub.CirclePublicMessage()
}
