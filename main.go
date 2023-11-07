package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"time"
)

var current_iteration float64 = 0

func calcul_total_iteration(min, max, number_caraceters int) float64 {
	var total float64 = 0
	for i := min; i <= max; i++ {
		total += math.Pow(float64(number_caraceters), float64(i))
	}
	return total
}
func generate_combination_length_between(min, max int, characterSet []rune, words chan string) {
	total_iteration := calcul_total_iteration(min, max, len(characterSet))
	var wg sync.WaitGroup
	for length := min; length <= max; length++ {
		wg.Add(1)
		go generate_combination(length, characterSet, total_iteration, words, &wg)
	}
	wg.Wait()
	close(words)
}

func generate_combination(length int, characterSet []rune, total_iteration float64, words chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	generate_combination_helper(length, characterSet, "", total_iteration, words)
}

func generate_combination_helper(length int, characterSet []rune, currentCombination string, total_iteration float64, result chan string) {
	if length == 0 {
		result <- currentCombination
		current_iteration++
		updateProgress(current_iteration, total_iteration)
		return
	}

	for _, char := range characterSet {
		generate_combination_helper(length-1, characterSet, currentCombination+string(char), total_iteration, result)
	}
}

func generateMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
func task(words chan string, searched_result string, wg *sync.WaitGroup) {
	defer wg.Done()
	for word := range words {
		if generateMD5(word) == searched_result {
			println("find ", word)
			break
		}
	}
}

func dispalyWords(words chan string) {
	for word := range words {
		println(word)
	}
}

func updateProgress(current, total float64) {
	fmt.Printf("\r \033[33m tested %.0f   progress:  %f%s  \033[0m",current, (current/total)*100, "%")
}
func main() {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789@"
	//characters := "AB"
	characterSet := []rune(characters)
	words := make(chan string)
	var wg sync.WaitGroup

	caracter_length_min:=4
	caracter_length_max:=5
	number_threads := 2
	searched_md5 := "16bb5086b1b21b2e86adbfefacfb7fbd"

	start_time:=time.Now()
	go generate_combination_length_between(caracter_length_min, caracter_length_max, characterSet, words)
	wg.Add(1)
	for i := 0; i < number_threads; i++ {
		go task(words, searched_md5, &wg)
	}
	wg.Wait()
	end_time:=time.Now()
	fmt.Printf("time taken is %v",end_time.Sub(start_time))

}
