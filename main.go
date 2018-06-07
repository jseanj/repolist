package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

const API_KEY string = "your api key"

var ids map[int64]int8
var page = flag.Int("p", 0, "page")
var show = flag.Bool("show", false, "show history")

func main() {
	flag.Parse()
	ids = make(map[int64]int8)
	if exists("./ids.data") {
		input, err := ioutil.ReadFile("./ids.data")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(input, &ids)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//fmt.Println(ids)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: API_KEY})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if *show && len(ids) > 0 {
		for id, _ := range ids {
			repo, _, _ := client.Repositories.GetByID(ctx, id)
			printRepo(*repo)
		}
		return
	}

	opts := &github.SearchOptions{}
	opts.Sort = "updated"
	opts.Order = "desc"
	opts.PerPage = 20
	//for i := 0; i < 1; i++ {
	opts.Page = *page
	repos, _, _ := client.Search.Repositories(ctx, "stars:>100 language:go", opts)
	printRepos(repos.Repositories)
	//}

	rate, _, _ := client.RateLimits(ctx)
	fmt.Printf("Rate Limit: %d/%d\n", rate.Core.Remaining, rate.Core.Limit)

	//fmt.Println(ids)
	output, err := json.Marshal(ids)
	if err != nil {
		fmt.Println(err)
		return
	}
	ioutil.WriteFile("./ids.data", output, 0644)
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func printRepos(repositories []github.Repository) {
	for _, repo := range repositories {
		id := *repo.ID
		if _, ok := ids[id]; ok {
			continue
		}
		printRepo(repo)
		ids[id] = 1
	}
}

func printRepo(repo github.Repository) {
	fmt.Print(color.YellowString("%s", *repo.HTMLURL), "  ", color.GreenString("%d", *repo.StargazersCount))
	if repo.Description != nil {
		fmt.Print(" ", color.WhiteString("%s", *repo.Description))
	}
	if repo.Homepage != nil {
		fmt.Println(" ", color.RedString("%s", *repo.Homepage))
	} else {
		fmt.Print("\n")
	}
}
