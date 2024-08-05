package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {

	var language_list []map[string]string

	language_list = []map[string]string{
		{"en": "English"},
		{"tr": "Turkish"},
		{"ar": "Arabic"},
		{"zh-CN": "Chinese (Simplified)"},
		{"zh-TW": "Chinese (Traditional)"},
		{"cs": "Czech"},
		{"da": "Danish"},
		{"nl": "Dutch"},
		{"fi": "Finnish"},
		{"de": "German"},
		{"fr": "French"},
		{"he": "Hebrew"},
		{"hi": "Hindi"},
		{"id": "Indonesian"},
		{"it": "Italian"},
		{"ja": "Japanese"},
		{"ko": "Korean"},
		{"ms": "Malay"},
		{"no": "Norwegian"},
		{"fa": "Persian"},
		{"pl": "Polish"},
		{"pt": "Portuguese"},
		{"ru": "Russian"},
		{"es": "Spanish"},
		{"sv": "Swedish"},
		{"si": "Sinhala"},
		{"ta": "Tamil"},
	}

	for _, v := range language_list {
		for k := range v {
			file_path := fmt.Sprintf("files/%s.json", k)
			file, err := os.Create(file_path)
			if err != nil {
				panic(err)
			}

			url := fmt.Sprintf("https://localeasy.dev/generate?sheet=1PLdjhmMotS60PMlds-VHFMVdoWW9e-Z_i1yKqDlTj_Q&format=json&locale=%s&nest-by=.&variant=web", k)
			fmt.Println(url)

			res, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			body, err := io.ReadAll(io.Reader(res.Body))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Define a struct to store the JSON data
			var data map[string]string

			// Unmarshal the JSON response body into the struct
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("Json Unmarshal Error:", err)
				//fmt.Println("Json Unmarshal Error:", string(body))
				return
			}

			// Print the JSON data (optional)
			fmt.Println("Response Data:", data)

			newMap := make(map[string]string)

			for key, value := range data {
				newkey := strings.ReplaceAll(key, " ", "_")
				newkey = strings.ToLower(newkey)
				if strings.Contains(newkey, "/") {
					newkey = strings.ReplaceAll(newkey, "/", "_or_")
				}

				if strings.Contains(newkey, "?") {
					newkey = strings.ReplaceAll(newkey, "?", "_or_")
				}

				if strings.Contains(newkey, ".") {
					newkey = strings.ReplaceAll(newkey, ".", "")
				}

				if strings.Contains(newkey, "\\") {
					newkey = strings.ReplaceAll(newkey, "\\", "")
				}

				if strings.Contains(newkey, "\u0026") {
					newkey = strings.ReplaceAll(newkey, "\u0026", "")
				}

				if strings.Contains(newkey, ":") {
					newkey = strings.ReplaceAll(newkey, ":", "")
				}

				if strings.Contains(newkey, "'") {
					newkey = strings.ReplaceAll(newkey, "'", "")
				}

				if strings.Contains(newkey, ",") {
					newkey = strings.ReplaceAll(newkey, ",", "")
				}

				if strings.Contains(newkey, "!") {
					newkey = strings.ReplaceAll(newkey, "!", "")
				}
				if strings.Contains(newkey, "-") {
					newkey = strings.ReplaceAll(newkey, "-", "_")
				}

				if strings.Contains(newkey, "\n") {
					newkey = strings.ReplaceAll(newkey, "\n", "")
				}

				if strings.Contains(newkey, "_{") {
					re := regexp.MustCompile(`_{.*}`)
					nk := re.ReplaceAllString(newkey, "")
					last := len(nk)
					fmt.Println(string(nk))
					if string(nk[last-1]) == "_" {
						newkey = string(nk[0 : last-1])
					}
					newkey = nk
				}

				if strings.Contains(newkey, "{") {
					re := regexp.MustCompile(`{.*}`)
					nk := re.ReplaceAllString(newkey, "")
					last := len(nk)
					fmt.Println(string(nk))
					if string(nk[last-1]) == "_" {
						newkey = string(nk[0 : last-1])
					}
					newkey = nk
				}

				newMap[newkey] = value
				//fmt.Println(newkey)
			}

			data = newMap

			outputContent, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshaling JSON data:", err)
				return
			}

			defer file.Close()

			// Encode and write the JSON data to the file
			// encoder := json.NewEncoder(file)
			// err = encoder.Encode(data)
			// if err != nil {
			// 	fmt.Println("Error:", err)
			// 	return
			// }

			nerr := os.WriteFile(file_path, outputContent, fs.FileMode(0644))
			if nerr != nil {
				fmt.Println("Error writing output file:", nerr)
				return
			}

			fmt.Println("Response saved as response.json")
		}
	}
}
