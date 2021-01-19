package main

import (
	"fmt"
	"github.com/sclevine/agouti"
	"log"
	"time"
)

const (
	CircleCIURL    = "https://circleci.com/login/"
	GithubUserName = ""
	GithubPassWord = ""
	Organization   = ""
	Repository     = ""
)

func main() {
	driver := agouti.ChromeDriver()

	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("Failed to open page: %v", err)
	}

	if err := page.Navigate(CircleCIURL); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	// Github Login
	page.FindByID("login_field").Fill(GithubUserName)
	page.FindByID("password").Fill(GithubPassWord)
	page.FindByButton("Sign in").Click()

	// Move to environment variables setting page
	page.Navigate(fmt.Sprintf("https://app.circleci.com/settings/project/github/%s/%s/environment-variables?return-to=https%%3A%%2F%%2Fapp.circleci.com%%2Fpipelines%%2Fgithub%%2F%s%%2F%s", Organization, Repository, Organization, Repository))

	if err := setVariables(page); err != nil {
		log.Fatalf("Failed to setVariables: %v", err)
	}

	log.Println("Success!")
}

func setVariables(page *agouti.Page) error {
	// 設定したい環境変数をKey, Value形式で記述する
	vars := map[string]string{
		"HOGE": "hoge",
		"PUGE": "puge",
		"FUGA": "fuga",
	}

	for k, v := range vars {
		time.Sleep(time.Second * 3) // DOMの生成を待つ
		page.FindByButton("Add Environment Variable").Click()
		page.FindByID("name").Fill(k)
		page.FindByID("value").Fill(v)
		err := page.Find("body > reach-portal > div:nth-child(3) > div > div > div > div.css-1m79yrn > form > div.css-1s3dnxk > button.css-1ywpecc.e1jkxlkv0").Click() // ここ不安定
		if err != nil {
			log.Printf("Failed to Click: %v", err)
			return err
		}
		fmt.Printf("登録完了: %s\n", k)
	}
	time.Sleep(time.Second * 1) // 最後の環境変数が登録される前に処理が終了するのを防ぐため
	return nil
}