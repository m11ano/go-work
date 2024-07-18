package main

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"hw3/user"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	r = regexp.MustCompile("@")
)

// вам надо написать более быструю оптимальную этой функции

func FastSearch(out io.Writer) {
	/*
		!!! !!! !!!
		обратите внимание - в задании обязательно нужен отчет
		делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировалм их
		так же обратите внимание на команду с параметром -http
		перечитайте еще раз задание
		!!! !!! !!!
		1) ReadAll - 10.5 MB - много
		2) strings.Split 2.2 MB - можно избежать
		3) json.Unmarshal 5.8 MB - можно оптимизировать через кодогенерацию easyjson
		4) regexp.MatchString - можно оптимизировать
	*/

	_, _ = fmt.Fprintln(out, "found users:")

	browsers := map[string]struct{}{}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	i := -1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		i++
		userItem := &user.JsonUser{}
		err := easyjson.Unmarshal(scanner.Bytes(), userItem)
		if err != nil {
			continue
		}

		isIE := false
		isAndroid := false

		for _, browser := range userItem.Browsers {
			anyTargetBrowser := false
			if strings.Contains(browser, "Android") {
				isAndroid = true
				anyTargetBrowser = true
			}

			if strings.Contains(browser, "MSIE") {
				isIE = true
				anyTargetBrowser = true
			}

			if anyTargetBrowser {
				browsers[browser] = struct{}{}
			}
		}

		if !(isIE && isAndroid) {
			continue
		}

		email := r.ReplaceAllString(userItem.Email, " [at] ")
		_, _ = fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, userItem.Name, email))
	}

	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
}
