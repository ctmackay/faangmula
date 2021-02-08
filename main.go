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
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var counter = make(map[string]int)
var companySlicer = make(map[string][]string)
var masterCompanyList []string
var words = readIntoArray("usa2.txt")
var catchPhrases = readIntoArray("catchphrases.txt")

func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}

func readIntoArray(src string) []string {
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

		// pop companyName from array
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

func initializeThings() (map[string]int, map[string][]string) {
	fmt.Println("hello")

	var data []string

	companyMap := make(map[string]int)

	file, _ := ioutil.ReadFile("company_list.json")
	json.Unmarshal(file, &data)
	// fmt.Println(data)

	// generate company map
	for _, companyName := range data {

		var firstChar = strings.ToUpper(string(companyName[0]))
		// fmt.Println(index, companyName, firstChar)

		masterCompanyList = append(masterCompanyList, companyName)
		companyMap[companyName] = 1

		if val, ok := counter[firstChar]; ok {

			counter[firstChar] = val + 1
			companySlicer[firstChar] = append(companySlicer[firstChar], companyName)

		} else {
			counter[firstChar] = 1
			companySlicer[firstChar] = []string{companyName}

		}

	}

	// iterate over company slicer and shuffle the slices

	for k := range companySlicer {
		companySlicer[k] = shuffle(companySlicer[k])
	}

	// fmt.Println(counter)
	return counter, companySlicer

}

func getRandDictWord() string {

	randInt := generateRandomInteger(0, len(words))

	randomWord := words[randInt]

	return randomWord
}

func getRandomAcronym() string {
	fmt.Println(companySlicer)

	shuffledArray := shuffle(masterCompanyList)

	numChars := generateRandomInteger(4, 10)
	fmt.Println("building a ", numChars, " character long string")

	var acronym = generateRandomAcronymFromCompanyList(numChars, shuffledArray)

	return acronym
}

func getRandomCatchphrase(acronym string) string {

	randInt := generateRandomInteger(0, len(catchPhrases)-1)
	fmt.Println(randInt)
	template := catchPhrases[randInt]
	catchPhrase := strings.ReplaceAll(template, "ACRONYM_REPLACE", acronym)

	fmt.Println("\n" + catchPhrase)
	return catchPhrase
}

func getCompaniesByAcronym(acronym string) []string {

	var _counter map[string]int = make(map[string]int)
	var _companySlicer map[string][]string = make(map[string][]string)

	// Copy from the original map to the target map
	for key, value := range companySlicer {
		_companySlicer[key] = value
	}
	// Copy from the original map to the target map
	for key, value := range counter {
		_counter[key] = value
	}

	var companies []string
	for _, chr := range acronym {

		var key string = strings.ToUpper(string(chr))

		var companyThatBeginsWithLetter = ""
		val := _counter[key]
		if val > 0 {
			_counter[key]--
			companyThatBeginsWithLetter, _companySlicer[key] = _companySlicer[key][0], _companySlicer[key][1:]

		} else {
			fmt.Println("Out of letters for this character. need to try a new word.")
			break
		}

		companies = append(companies, companyThatBeginsWithLetter)

	}
	return companies
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func routeGetCompaniesByAcronym(c *gin.Context) {
	var ac string = c.Param("acronym")
	var companyList = getCompaniesByAcronym(ac)
	c.JSON(http.StatusOK, gin.H{"acronym": ac, "companies": companyList, "catchphrase": getRandomCatchphrase(ac)})
}

func routeGetRandomDictWord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"companies": getRandDictWord()})
}

func routeGetRandomAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"companies": getRandomAcronym()})
}

func routeGetRandomCatchphrase(c *gin.Context) {
	var ac string = c.Param("acronym")
	c.JSON(http.StatusOK, gin.H{"companies": getRandomCatchphrase(ac)})
}

func routeGetAllCompanyList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"companies": masterCompanyList})
}

func routeGetRandomSet(c *gin.Context) {
	acronym := getRandomAcronym()
	catchPhrase := getRandomCatchphrase(acronym)
	companies := getCompaniesByAcronym(acronym)
	c.JSON(http.StatusOK, gin.H{"companies": companies, "catchphrase": catchPhrase, "acronym": acronym})
}

func routeGetRandomSetDictWord(c *gin.Context) {
	acronym := getRandDictWord()
	catchPhrase := getRandomCatchphrase(acronym)
	companies := getCompaniesByAcronym(acronym)
	c.JSON(http.StatusOK, gin.H{"companies": companies, "catchphrase": catchPhrase, "acronym": acronym})
}

func main() {
	initializeThings()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.GET("/getCompaniesByAcronym/:acronym", routeGetCompaniesByAcronym)
	r.GET("/getAllCompanies", routeGetAllCompanyList)
	r.GET("/getRandomDictWord", routeGetRandomDictWord)
	r.GET("/getRandomCatchPhrase/:acronym", routeGetRandomCatchphrase)
	r.GET("/getRandomAcronym", routeGetRandomAcronym)
	r.GET("/getRandomSet", routeGetRandomSet)
	r.GET("/getRandomDictSet", routeGetRandomSetDictWord)

	port := ":10000"
	r.Run(port)
}
