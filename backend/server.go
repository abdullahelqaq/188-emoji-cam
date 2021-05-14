package main

import (
	"fmt"
	"log"
	"net/http"
	"bufio"
	"os"
	"strings"
	"sort"
	"encoding/json"
)

var words []string

type ByLen []string
 
func (a ByLen) Len() int {
   return len(a)
}
 
func (a ByLen) Less(i, j int) bool {
   return len(a[i]) < len(a[j])
}
 
func (a ByLen) Swap(i, j int) {
   a[i], a[j] = a[j], a[i]
}

func loadWords(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("Failed to open file!")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	file.Close()

	/*
	for _, each_ln := range words {
		fmt.Println(each_ln)
	}
	*/
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {

	loadWords("words_60k.txt")
	//loadWords("words_alpha.txt")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		payload := strings.Split(r.URL.Path[1:], ",")
		entry, next_letters := payload[0], payload[1]

		num_suggestions := 6
		if next_letters == "rtyu" {
			num_suggestions = 3
		}

		var all_options [][]string
		for i := 0; i < len(next_letters); i++ {
			potential_entry := entry + string(next_letters[i])
			var options []string
			for j := 0; j < len(words); j++ {
				if strings.HasPrefix(words[j], potential_entry) && len(words[j]) < 7 {
					options = append(options, words[j])
				}
			}
			sort.Sort(ByLen(options))
			all_options = append(all_options, options)
		}

		var final_options []string
		for len(final_options) < num_suggestions {
			all_empty := true
			for i := 0; i < len(all_options); i++ {

				// Put here for rtyu case /:
				if len(final_options) == num_suggestions {
					break
				}

				if len(all_options[i]) > 0 {
					all_empty = false
					final_options = append(final_options, all_options[i][0])
					all_options[i] = all_options[i][1:]
				}
			}
			if all_empty {
				break
			}
		}
		options_json, _ := json.Marshal(final_options)
		fmt.Fprintf(w, string(options_json))
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
