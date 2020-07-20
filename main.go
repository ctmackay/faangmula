package main

// basic logic
/*

grab list of companies.
grab their first letter.

create a scrambler and randomizer algorithm that picks 4-10 characters
from the available characters.

RNG 4-10,
eg 5:

datastructure we could use
a simple array with all first letters of each company shuffled.

a dictionary, and just random sampled keys.
then extract that from

display acronyms on the front page
FAANG


*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}

func readEnglishWords(src string) []string {
	var lines []string

	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())

		// fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines

}
func generateRandomAcronymFromCompanyList(numChars int, shuffledArray []string) string {
	var acronym string
	for i := 0; i < numChars; i++ {

		var company string

		// pop element from array
		company, shuffledArray = shuffledArray[0], shuffledArray[1:]
		fmt.Println(company)
		var firstChar = string(company[0])

		acronym += firstChar

	}
	fmt.Print(acronym)
	return acronym
}

func generateRandomInteger(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	numChars := rand.Intn(max-min+1) + min
	return numChars
}

func main() {
	fmt.Println("hello")

	var data []string

	counter := make(map[string]int)
	companyMap := make(map[string]int)
	companySlicer := make(map[string][]string)

	file, _ := ioutil.ReadFile("company_list.json")
	json.Unmarshal(file, &data)
	// fmt.Println(data)

	// generate company map
	for _, element := range data {

		var firstChar = strings.ToUpper(string(element[0]))
		// fmt.Println(index, element, firstChar)

		companyMap[element] = 1

		if val, ok := counter[firstChar]; ok {

			counter[firstChar] = val + 1
			companySlicer[firstChar] = append(companySlicer[firstChar], element)

		} else {
			counter[firstChar] = 1
		}

	}

	// iterate over company slicer and shuffle the slices

	for k := range companySlicer {
		companySlicer[k] = shuffle(companySlicer[k])
	}

	fmt.Println(companySlicer)

	shuffledArray := shuffle(data)

	numChars := generateRandomInteger(4, 10)
	fmt.Println("building a ", numChars, " character long string")

	var acronym = generateRandomAcronymFromCompanyList(numChars, shuffledArray)
	catchPhrase := "its fucked up working for " + acronym
	fmt.Println("\n" + catchPhrase)

	words := readEnglishWords("usa2.txt")

	randInt := generateRandomInteger(0, len(words))

	randomWord := words[randInt]

	fmt.Println(randomWord)

	fmt.Println(counter)

	for i, chr := range randomWord {

		var key string = strings.ToUpper(string(chr))

		var companyThatBeginsWithLetter = ""
		val := counter[key]
		if val > 0 {
			counter[key]--
			companyThatBeginsWithLetter, companySlicer[key] = companySlicer[key][0], companySlicer[key][1:]

		} else {
			fmt.Println("Out of letters for this character. need to try a new word.")
			break
		}

		fmt.Printf("index: %2v  chr: %c remaining: %d %s\n", i, chr, val, companyThatBeginsWithLetter)

	}

}
