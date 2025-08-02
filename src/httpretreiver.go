package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"bufio"
	"github.com/schollz/progressbar/v3"
)

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

func filter(regex string, data2 string) ([]string, []string) {
    r := regexp.MustCompile(regex)
	  matches := r.FindAllString(data2, -1)
    
    codes := []string{}
    file_names := []string{}
	  for i, v := range matches {
			// get the file names
			file_name := v[strings.Index(v, "\">")+2:strings.Index(v, "</a>")]
      file_names = append(file_names, file_name)
			// get the ids
      code := v[strings.Index(v, ";query=")+7:strings.Index(v, "\" title")]
      codes = append(codes, code)
			// get the song titles
			titles := v[strings.Index(v, "title=\"")+7:strings.Index(v, "\">")]
			
		fmt.Println(i, "-", file_name, "| Song Title: \""+titles+"\"")
	  }
    return codes, file_names
}

func retrieve(query string, page_no int) ([]string, []string) {
  var query_codes []string
  var query_names []string
	response, err := http.Get("https://modarchive.org/index.php?request=search&query="+query+"&submit=Find&search_type=filename_or_songtitle&page="+strconv.Itoa(page_no)+"#mods")
	if err != nil {
		fmt.Printf("The http request failed with error: %s\n", err)
	} else {
		data, _ := io.ReadAll(response.Body)
    
    if strings.Contains(string(data), "Sorry, no results were found.") {
      fmt.Println("Sorry, no results were found.")
			return nil, nil
    } else if strings.Contains(string(data), "Your search query was too short (minimum 3 alphanumeric characters)") {
			fmt.Println("Your search query was too short (minimum 3 alphanumeric characters)")
			return nil, nil
		} else {
      // extract the hyperlinks
      var data2 string = string(data)
      query_codes, query_names = filter(`<a class="standard-link".*`, data2)
    }
		
    //data2 := strings.Contains(string(data), "<a class=\"standard-link\"")
		//fmt.Println(string(data2))
	}
  return query_codes, query_names
}

func lookupFileName(query_id string) {
  response, err := http.Get("https://modarchive.org/index.php?request=view_by_moduleid&query="+query_id)
  if err != nil {
    return
  }
  fmt.Println(response)
}

func helperFunc(query *string, selection *int, codes *[]string, names *[]string) {
  var query2 string
	reader := bufio.NewReader(os.Stdin)

	for {
  	fmt.Print("Search by name for: ")
		//_, err := fmt.Scan(&query2)
		answer, _ := reader.ReadString('\n')
		
		if answer == "\n" {
			query2 = ""
		} else {
			answer = strings.ReplaceAll(answer, " ", "+")
			query2 = answer[:len(answer)-1]
		}
		//fmt.Println(answer, len(answer))
		if query2 == "" {
			fmt.Println("Input cannot be empty")
			continue
		}
    //if err != nil {
		//	fmt.Println("Error reading query", err)
		//	continue
		//} 

		break
	}
	fmt.Println("Your search is : ", query2)

	clearConsole()
  
  var page int;
  page = 1;
  *query = query2
  
  for {
    fmt.Println("===== search results in page "+strconv.Itoa(page)+" =====")
    *codes, *names = retrieve(*query, page)
    if (*codes == nil && *names == nil) { break }
    var page_input string

    fmt.Print("\n\nPage number (type \"quit\" to quit): ")
    fmt.Scan(&page_input)
    page, _ = strconv.Atoi(page_input)
    if strings.Contains(page_input, "quit") {
      break
    } else {
      fmt.Println("invalid command")
    }
		clearConsole()
    fmt.Println("Looking in page ", page)
  }
  //fmt.Println(codes)
  for {
		if (*codes == nil && *names == nil) { break }
    fmt.Print("Enter chosen index of file for current page: ")
    fmt.Scanf("%d", &(*selection))

    if *selection > len((*codes))-1 || *selection < 0 {
      fmt.Println("Error: Invalid index")
    } else {
      break
    }
  }
}

// func main() {
// 	var query string
//   var selection int
//   var codes []string
//   var names []string

//   for {
// 		helperFunc(&query, &selection, &codes, &names)
// 		if (codes != nil && names != nil) { break; }
// 	}
// 	downloadFile(names[selection], "https://api.modarchive.org/downloads.php?moduleid="+codes[selection])
// }

func downloadFile(filepath string, url string) (err error) {
	dirpath := "downloads/"

	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err := os.MkdirAll(dirpath, os.ModePerm)

		// if for some reason that fails too
		if err != nil {
			fmt.Println("error making directory: ", err)
			return err
		}
		fmt.Println("Directory created")
	} else {
		fmt.Println("Directory already exists")
	}
  // Create the file
  out, err := os.Create(dirpath+filepath)

  if err != nil  {
    return err
  }
  defer out.Close()

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
	
	bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading: "+filepath)
  // Writer the body to file
  _, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
  if err != nil  {
    return err
  }

  return nil
}
