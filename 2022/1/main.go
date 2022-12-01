package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("./data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	elfMap := []int{0}
	currentIndex := 0

	var tempText string
	var tempInt int
	for scanner.Scan() {
		tempText = scanner.Text()
		if tempText == "" {
			currentIndex++
			elfMap = append(elfMap, 0)
			continue
		}
		tempInt, err = strconv.Atoi(tempText)
		if err != nil {
			log.Fatal(err)
		}
		elfMap[currentIndex] += tempInt
	}
	var thirdI int
	var thirdV int
	var secondI int
	var secondV int
	var firstI int
	var firstV int
	for i, v := range elfMap {
		switch {
		case v > firstV:
			thirdI = secondI
			thirdV = elfMap[thirdI]
			secondI = firstI
			secondV = elfMap[secondI]
			firstI = i
			firstV = elfMap[firstI]
		case v > secondV:
			thirdI = secondI
			thirdV = elfMap[thirdI]
			secondI = i
			secondV = elfMap[secondI]
		case v > thirdV:
			thirdI = i
			thirdV = elfMap[thirdI]
		}
		log.Println("elf ", i, " : ", v)
	}
	log.Println("FIRST - elf ", firstI, " : ", firstV)
	log.Println("SECOND - elf ", secondI, " : ", secondV)
	log.Println("THIRD - elf ", thirdI, " : ", thirdV)
	log.Println("TOP 3 SUM - ", firstV+secondV+thirdV)
}
