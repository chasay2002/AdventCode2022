package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
As you walk, the Elf shows you a small bag and some cubes which are either red, green, or blue. Each time you play this game, he will hide a secret number of cubes of each color in the bag, and your goal is to figure out information about the number of cubes.

To get information, once a bag has been loaded with cubes, the Elf will reach into the bag, grab a handful of random cubes, show them to you, and then put them back in the bag. He'll do this a few times per game.

You play several games and record the information from each game (your puzzle input). Each game is listed with its ID number (like the 11 in Game 11: ...) followed by a semicolon-separated list of subsets of cubes that were revealed from the bag (like 3 red, 5 green, 4 blue).

For example, the record of a few games might look like this:

Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
In game 1, three sets of cubes are revealed from the bag (and then put back again). The first set is 3 blue cubes and 4 red cubes; the second set is 1 red cube, 2 green cubes, and 6 blue cubes; the third set is only 2 green cubes.

The Elf would first like to know which games would have been possible if the bag contained only 12 red cubes, 13 green cubes, and 14 blue cubes?

In the example above, games 1, 2, and 5 would have been possible if the bag had been loaded with that configuration. However, game 3 would have been impossible because at one point the Elf showed you 20 red cubes at once; similarly, game 4 would also have been impossible because the Elf showed you 15 blue cubes at once. If you add up the IDs of the games that would have been possible, you get 8.

Determine which games would have been possible if the bag had been loaded with only 12 red cubes, 13 green cubes, and 14 blue cubes. What is the sum of the IDs of those games?
*/

/*
As you continue your walk, the Elf poses a second question: in each game you played, what is the fewest number of cubes of each color that could have been in the bag to make the game possible?

Again consider the example games from earlier:

Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
- In game 1, the game could have been played with as few as 4 red, 2 green, and 6 blue cubes. If any color had even one fewer cube, the game would have been impossible.
- Game 2 could have been played with a minimum of 1 red, 3 green, and 4 blue cubes.
- Game 3 must have been played with at least 20 red, 13 green, and 6 blue cubes.
- Game 4 required at least 14 red, 3 green, and 15 blue cubes.
- Game 5 needed no fewer than 6 red, 3 green, and 2 blue cubes in the bag.
The power of a set of cubes is equal to the numbers of red, green, and blue cubes multiplied together. The power of the minimum set of cubes in game 1 is 48. In games 2-5 it was 12, 1560, 630, and 36, respectively. Adding up these five powers produces the sum 2286.

For each game, find the minimum set of cubes that must have been present. What is the sum of the power of these sets?
*/

type Result struct {
	gameID   int
	minBlue  int
	minGreen int
	minRed   int
	minTotal int
}

type Game struct {
	ID   int
	Sets []Set
}

type Set struct {
	Blue  int
	Green int
	Red   int
}

var id = regexp.MustCompile(`Game (?P<id>\d+)`)
var groups = regexp.MustCompile(`(?P<count>\d+) (?P<color>blue|red|green)`)

var input = `Game 1: 2 green, 12 blue; 6 red, 6 blue; 8 blue, 5 green, 5 red; 5 green, 13 blue; 3 green, 7 red, 10 blue; 13 blue, 8 red
Game 2: 1 green, 7 red; 1 green, 9 red, 3 blue; 4 blue, 5 red
Game 3: 2 red, 2 blue, 6 green; 1 blue, 2 red, 2 green; 3 blue, 3 green
Game 4: 8 green, 16 red, 7 blue; 1 red, 7 blue, 12 green; 8 green, 14 red, 1 blue; 6 blue, 9 green, 12 red; 9 red, 2 green; 8 red, 7 blue, 17 green
Game 5: 2 red, 7 green; 2 red, 1 green, 4 blue; 4 blue, 7 green, 5 red; 8 red, 2 blue, 5 green
Game 6: 20 red, 4 green, 3 blue; 12 green, 1 blue, 9 red; 8 green, 15 red, 1 blue; 17 red, 5 blue, 7 green
Game 7: 2 red, 2 blue, 5 green; 3 blue, 1 green, 2 red; 2 green, 2 blue; 2 green, 1 red, 1 blue; 1 blue, 6 green; 1 green, 2 red
Game 8: 18 red, 1 green, 3 blue; 7 blue, 3 green, 18 red; 2 green
Game 9: 13 blue; 13 blue; 3 red, 13 blue; 1 red, 8 blue; 1 green, 8 blue, 1 red
Game 10: 15 blue, 1 green, 12 red; 13 red, 13 green, 5 blue; 2 blue, 3 red, 7 green
Game 11: 2 blue, 17 red, 5 green; 6 red, 2 green, 12 blue; 7 green, 5 blue, 14 red
Game 12: 12 blue, 10 green, 14 red; 7 red, 3 green; 10 blue, 2 red; 4 red, 5 blue, 7 green
Game 13: 1 red, 13 blue, 14 green; 1 green, 1 red, 7 blue; 5 red, 15 blue, 14 green
Game 14: 2 blue, 1 red, 9 green; 4 blue, 3 red; 2 red, 7 green; 2 red, 4 green
Game 15: 5 green, 7 blue, 18 red; 10 blue, 4 green, 4 red; 13 blue, 16 red; 18 red, 4 green; 1 blue, 6 green, 17 red; 9 blue, 2 green, 15 red
Game 16: 15 green, 3 red, 4 blue; 3 blue, 16 green; 5 green, 2 blue, 1 red; 6 blue, 3 red, 12 green; 2 red, 3 green
Game 17: 10 green, 4 blue, 11 red; 6 green, 13 red, 6 blue; 2 blue, 5 green, 2 red
Game 18: 1 red, 13 blue; 4 green, 5 red; 3 green, 1 red, 12 blue; 1 green, 15 blue; 7 red, 1 green, 15 blue
Game 19: 2 blue, 6 green, 4 red; 4 green, 15 red; 3 blue, 8 red, 7 green
Game 20: 7 blue, 7 red, 10 green; 5 red, 1 green, 9 blue; 1 blue, 6 red, 12 green; 1 red, 11 green; 6 green, 9 blue, 4 red
Game 21: 4 red, 16 green, 10 blue; 2 blue, 14 green; 6 green, 4 blue; 7 blue, 6 red; 6 red, 10 green, 6 blue
Game 22: 5 green, 2 red; 5 green, 3 red, 1 blue; 2 red, 1 green, 7 blue; 3 red, 11 green, 4 blue; 6 blue, 10 green, 2 red; 5 blue, 2 green, 2 red
Game 23: 3 blue, 3 red; 6 blue, 6 green, 2 red; 3 green, 15 blue, 2 red; 2 green, 1 blue, 2 red; 5 green, 5 blue; 5 green, 10 blue, 7 red
Game 24: 1 blue, 19 green; 18 green, 1 red, 2 blue; 5 green; 2 blue, 3 green, 1 red; 1 red, 5 green; 9 green, 1 red
Game 25: 1 green, 1 red, 7 blue; 2 green, 10 blue, 1 red; 8 red, 1 blue, 1 green
Game 26: 7 blue, 2 green, 2 red; 8 blue, 6 green, 9 red; 5 red, 11 green, 2 blue; 7 blue, 5 red, 9 green; 10 blue, 12 green, 4 red
Game 27: 9 blue, 1 red; 15 blue, 2 red; 2 red, 11 blue; 15 blue, 4 red, 1 green; 4 red, 1 green
Game 28: 9 green, 7 blue; 8 green, 3 red, 2 blue; 3 red, 11 blue, 7 green; 3 red, 3 blue; 9 green, 2 red
Game 29: 2 green, 3 blue, 1 red; 13 blue, 2 green, 8 red; 3 green, 12 red, 9 blue
Game 30: 9 red, 18 green, 6 blue; 1 green, 3 blue, 5 red; 6 blue, 12 green, 3 red
Game 31: 6 green, 1 red; 1 blue, 4 red, 10 green; 11 green, 1 blue, 4 red; 2 green; 1 red, 1 green; 8 green
Game 32: 7 green, 8 blue, 1 red; 4 red, 5 blue, 8 green; 5 green, 2 blue; 2 red, 3 green, 2 blue
Game 33: 1 red, 1 blue; 16 red, 1 green, 2 blue; 1 blue, 13 red; 1 blue, 1 green, 17 red; 6 red, 1 blue, 1 green; 10 red, 2 blue
Game 34: 10 red, 13 green; 18 blue, 6 green; 5 blue, 1 red, 11 green; 8 green, 1 red, 9 blue; 11 blue, 10 red, 12 green
Game 35: 4 red, 2 green; 4 red, 3 blue, 7 green; 3 blue, 2 green, 7 red
Game 36: 1 blue, 8 red, 9 green; 2 green, 9 red, 2 blue; 8 blue, 10 green, 12 red; 9 blue, 9 green, 10 red; 1 red, 9 green, 5 blue
Game 37: 7 blue, 2 green, 14 red; 8 green, 8 red, 6 blue; 2 green, 12 red, 1 blue; 5 blue, 6 green, 5 red; 9 blue, 2 green, 1 red
Game 38: 11 red, 2 green; 15 red, 1 blue, 7 green; 3 red, 3 green, 1 blue; 2 blue, 11 red, 9 green; 2 blue, 3 green; 5 green, 8 red, 1 blue
Game 39: 4 blue, 13 red; 16 red, 3 blue, 2 green; 1 red, 1 blue
Game 40: 7 blue, 2 green; 2 green, 2 blue; 1 red, 1 blue; 1 red, 4 green, 15 blue
Game 41: 3 red, 5 blue, 5 green; 1 green, 4 blue, 1 red; 9 blue, 3 green, 1 red; 3 green, 9 blue, 1 red
Game 42: 8 green, 2 blue, 4 red; 10 green, 2 red, 1 blue; 4 red, 10 green; 5 green, 3 blue, 3 red; 5 green, 3 blue, 5 red; 3 red, 10 green, 5 blue
Game 43: 14 red, 8 green, 3 blue; 9 red, 2 blue, 9 green; 6 red, 4 green, 4 blue; 8 green; 12 blue, 10 red, 9 green
Game 44: 6 green, 1 red; 3 red, 11 green, 2 blue; 2 green, 2 red, 3 blue; 1 red, 15 green, 2 blue
Game 45: 12 green, 16 red, 6 blue; 15 blue, 3 red, 5 green; 2 blue, 10 green; 15 red, 15 blue; 6 green, 5 blue, 16 red; 16 green, 6 blue
Game 46: 1 green, 3 red, 18 blue; 3 red, 8 green, 10 blue; 3 green, 12 blue; 2 red, 10 green, 17 blue; 12 green, 9 blue; 3 green, 6 blue
Game 47: 10 green, 7 red, 6 blue; 12 green, 2 blue, 5 red; 7 green, 10 blue, 3 red; 11 green, 7 red; 5 blue, 8 green
Game 48: 1 green, 2 red; 2 blue, 1 green, 6 red; 5 red, 1 blue; 2 blue, 1 green, 1 red; 9 red
Game 49: 6 green, 4 blue, 1 red; 6 blue, 2 red, 1 green; 2 red, 5 blue, 4 green; 3 green, 6 blue, 1 red; 6 green, 12 blue, 1 red
Game 50: 2 blue, 4 red, 2 green; 6 red, 12 green, 4 blue; 1 blue, 6 red, 13 green
Game 51: 3 red, 1 blue; 2 red, 1 green, 2 blue; 3 blue, 2 red, 1 green; 1 red, 1 green, 1 blue
Game 52: 1 green, 12 blue; 2 red, 9 blue; 1 green, 8 blue, 2 red; 8 blue, 1 red; 1 red, 13 blue; 2 red, 2 blue, 1 green
Game 53: 7 blue, 16 green, 13 red; 9 red, 5 blue, 15 green; 4 green, 1 red; 15 green, 10 red, 2 blue; 11 green, 11 red, 2 blue
Game 54: 5 red, 20 green, 11 blue; 13 red, 4 blue; 13 green, 15 blue; 1 red, 8 blue, 11 green; 5 green, 17 blue, 4 red
Game 55: 5 red, 17 blue, 13 green; 17 red, 1 blue, 6 green; 3 red, 5 green; 8 green, 16 blue, 5 red
Game 56: 6 green, 5 blue; 7 green, 1 red, 1 blue; 7 green, 2 blue; 3 blue, 8 green; 5 blue
Game 57: 1 blue, 2 red, 3 green; 2 blue; 4 green, 2 red, 1 blue; 2 red, 1 green, 1 blue; 2 red, 2 blue, 1 green; 1 blue, 2 green
Game 58: 6 red, 7 green, 18 blue; 9 blue, 4 green, 8 red; 12 green, 5 blue, 7 red
Game 59: 2 red, 19 green, 1 blue; 10 green, 5 red; 2 blue, 10 green; 1 red, 3 blue; 14 green, 4 red
Game 60: 12 blue, 2 red, 11 green; 6 red, 1 blue, 11 green; 6 red, 11 blue, 7 green; 5 green, 4 red, 7 blue; 2 red, 14 green
Game 61: 7 red, 4 blue, 1 green; 11 blue, 8 red, 3 green; 10 blue, 7 red, 2 green; 2 green, 12 blue, 10 red
Game 62: 3 blue, 19 green, 1 red; 3 red, 7 blue, 14 green; 10 green, 1 blue
Game 63: 11 green, 1 blue, 15 red; 11 green, 16 red; 11 green, 8 blue, 12 red
Game 64: 8 red, 2 green, 3 blue; 1 green, 9 red, 7 blue; 11 red, 3 blue; 2 blue, 3 green, 9 red; 3 blue, 6 red
Game 65: 1 red, 4 blue; 4 blue, 6 red; 4 blue, 11 red; 2 green, 2 blue, 4 red; 1 blue, 9 red; 5 green, 7 red
Game 66: 8 blue, 11 green; 7 red, 6 green; 6 red, 7 blue, 9 green; 12 blue, 2 red, 6 green
Game 67: 4 green, 6 blue, 2 red; 4 blue, 2 red, 2 green; 4 green, 2 blue; 2 green, 4 blue, 2 red; 4 red, 4 green, 2 blue; 2 red, 7 green
Game 68: 1 red, 5 blue; 5 blue, 16 green, 2 red; 1 red, 12 green, 4 blue; 7 blue, 9 red, 6 green
Game 69: 6 blue, 6 green; 3 blue, 9 red; 3 green, 13 red, 6 blue
Game 70: 1 green, 5 red; 1 blue, 1 green, 15 red; 5 blue, 5 red; 10 red, 7 blue; 9 red, 7 blue; 4 red, 1 green, 3 blue
Game 71: 6 green, 11 blue, 6 red; 11 blue, 3 red, 4 green; 6 red, 10 blue
Game 72: 5 green, 3 blue, 6 red; 4 blue, 2 red; 3 red, 4 green, 1 blue; 18 blue, 4 green; 3 red, 1 blue, 1 green; 4 blue, 4 green, 1 red
Game 73: 8 green, 4 red, 13 blue; 2 green, 11 blue, 4 red; 2 blue, 12 red, 4 green; 2 blue, 1 red, 2 green; 11 blue, 8 red, 11 green
Game 74: 4 red, 8 blue, 3 green; 7 red, 7 blue, 7 green; 7 red, 2 blue, 1 green
Game 75: 1 blue, 1 green, 2 red; 6 red, 1 blue; 15 red, 13 blue; 19 red, 4 green, 4 blue; 12 blue, 3 red, 2 green; 12 red, 3 green, 1 blue
Game 76: 1 blue, 3 red, 12 green; 13 green, 1 red, 1 blue; 6 blue, 3 green, 2 red; 1 green, 1 blue
Game 77: 4 blue; 1 red, 1 green, 5 blue; 2 blue; 4 blue, 2 red, 3 green; 1 green, 5 blue; 1 red, 1 green, 6 blue
Game 78: 2 red, 4 green, 1 blue; 4 green, 1 blue, 6 red; 7 green, 1 blue
Game 79: 3 red, 14 blue, 1 green; 1 blue, 5 red; 4 green, 3 blue, 5 red; 7 red, 1 blue, 3 green; 7 blue, 4 green
Game 80: 12 red, 1 green, 11 blue; 2 green, 5 blue, 13 red; 15 red, 15 blue, 1 green; 1 red, 1 green, 14 blue
Game 81: 3 green, 2 red, 6 blue; 4 green, 11 blue; 4 green, 9 blue, 4 red; 10 blue, 2 red, 4 green
Game 82: 15 blue, 7 green; 4 green, 15 blue, 7 red; 4 green, 9 blue, 5 red; 4 green, 6 blue, 1 red
Game 83: 16 green, 10 red, 9 blue; 14 green, 3 blue, 7 red; 11 red, 7 blue; 10 blue, 15 red; 3 red, 4 blue, 10 green; 15 red, 3 blue
Game 84: 3 red, 14 blue; 1 red, 3 blue; 2 red, 4 green, 6 blue; 1 green, 11 blue, 2 red
Game 85: 15 green, 2 blue, 11 red; 4 red, 7 blue, 6 green; 3 blue, 14 green; 10 green, 10 blue, 7 red; 17 red, 3 green, 15 blue
Game 86: 8 red, 1 green, 14 blue; 6 red, 12 green, 9 blue; 13 blue, 9 green, 6 red
Game 87: 4 red, 3 blue; 4 blue, 4 green, 4 red; 1 green, 6 blue, 4 red; 2 green, 4 red, 3 blue; 4 blue, 4 green, 2 red; 1 green, 7 blue
Game 88: 6 blue, 4 red, 6 green; 1 red, 9 blue, 13 green; 3 red, 1 blue, 15 green; 3 red, 7 blue, 15 green; 1 green, 1 blue, 3 red; 4 green, 4 blue, 2 red
Game 89: 16 green, 9 blue, 1 red; 3 blue, 6 green, 5 red; 5 blue, 1 green
Game 90: 17 red, 1 blue; 5 red; 2 red, 1 blue, 1 green; 13 red; 11 red, 1 blue, 1 green
Game 91: 4 red, 1 blue; 4 green, 7 red; 13 green, 9 blue; 5 green, 9 blue; 12 red, 6 green, 6 blue
Game 92: 4 red, 1 blue; 2 red, 15 blue; 2 blue, 1 green, 8 red; 1 green, 6 blue, 8 red; 12 blue, 1 green, 8 red; 8 blue, 1 red
Game 93: 8 red, 9 green, 3 blue; 2 green, 1 red, 2 blue; 7 red, 5 blue
Game 94: 1 green, 2 blue, 10 red; 7 red, 3 blue; 1 blue, 6 red, 4 green
Game 95: 1 red, 7 blue, 2 green; 3 red, 14 blue, 2 green; 1 red; 1 red, 14 blue, 1 green; 4 blue, 10 red, 2 green; 9 blue, 7 red
Game 96: 1 green, 10 red, 1 blue; 4 red, 1 green, 2 blue; 3 blue, 8 red
Game 97: 4 red, 3 green; 1 blue, 2 green; 4 green, 4 red; 1 red, 1 blue; 1 green, 3 red; 1 green
Game 98: 7 green, 4 blue, 1 red; 2 red, 5 blue; 4 blue, 3 red, 7 green
Game 99: 7 blue, 13 green; 3 green, 5 red, 12 blue; 2 blue, 14 green, 8 red; 4 red, 6 blue, 2 green; 5 red, 9 green, 13 blue; 8 red, 8 blue, 5 green
Game 100: 8 green, 7 blue, 1 red; 10 blue, 2 green, 5 red; 12 blue, 1 green, 1 red; 9 green, 9 blue, 2 red; 1 blue, 5 red, 3 green`

var requirements = Set{
	Red:   12,
	Green: 13,
	Blue:  14,
}

func main() {
	in := input

	games := []Game{}
	for _, line := range strings.Split(strings.TrimSuffix(in, "\n"), "\n") {
		// fmt.Printf("%#v", ParseGame(line))
		games = append(games, ParseGame(line))
	}
	fmt.Printf("%#v\n", games)

	results := []Result{}
	for _, game := range games {
		results = append(results, GameToResult(game))
	}
	fmt.Printf("%#v\n", results)

	checkSum := 0
	rootSum := 0
	for _, result := range results {
		rootSum += result.minBlue * result.minGreen * result.minRed
		if result.minBlue > requirements.Blue ||
			result.minGreen > requirements.Green ||
			result.minRed > requirements.Red ||
			result.minTotal > requirements.Blue+requirements.Green+requirements.Red {
			continue
		}
		fmt.Printf("- %#v\n", result)
		checkSum += result.gameID
	}
	println(checkSum)
	println(rootSum)
}

func ParseGame(data string) Game {
	var err error

	game := Game{
		Sets: []Set{},
	}

	game.ID, err = strconv.Atoi(id.FindStringSubmatch(data)[1])
	if err != nil {
		panic(err)
	}

	sets := strings.Split(data, ";")
	fmt.Printf("%#v\n", sets)
	for _, setData := range sets {
		matches := groups.FindAllStringSubmatch(setData, -1)
		fmt.Printf("- %#v\n", matches)
		game.Sets = append(game.Sets, MatchesToSet(matches))
	}

	return game
}

func MatchesToSet(matches [][]string) Set {
	set := Set{}
	for _, match := range matches {
		fmt.Printf("-- %#v\n", match)
		if num, err := strconv.Atoi(match[1]); err == nil {
			switch match[2] {
			case "blue":
				set.Blue = num
			case "green":
				set.Green = num
			case "red":
				set.Red = num
			}
		}
	}

	return set
}

func GameToResult(game Game) Result {
	res := Result{
		gameID: game.ID,
	}

	for _, set := range game.Sets {
		fmt.Printf(">> %#v\n", set)
		if set.Blue > res.minBlue {
			res.minBlue = set.Blue
		}
		if set.Green > res.minGreen {
			res.minGreen = set.Green
		}
		if set.Red > res.minRed {
			res.minRed = set.Red
		}
		sum := set.Blue + set.Green + set.Red
		if sum > res.minTotal {
			res.minTotal = sum
		}
	}

	return res
}
