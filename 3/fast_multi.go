package main

//var (
//	r = regexp.MustCompile("@")
//)

// вам надо написать более быструю оптимальную этой функции
//
//	func FastSearch(out io.Writer) {
//		/*
//			!!! !!! !!!
//			обратите внимание - в задании обязательно нужен отчет
//			делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировалм их
//			так же обратите внимание на команду с параметром -http
//			перечитайте еще раз задание
//			!!! !!! !!!
//			1) ReadAll - 10.5 MB - много
//			2) strings.Split 2.2 MB - можно избежать
//			3) json.Unmarshal 5.8 MB - можно оптимизировать через кодогенерацию easyjson
//			4) regexp.MatchString - можно оптимизировать
//		*/
//
//		_, _ = fmt.Fprintln(out, "found users:")
//
//		browsers := map[string]struct{}{}
//
//		file, err := os.Open(filePath)
//		if err != nil {
//			panic(err)
//		}
//		defer file.Close()
//
//		i := -1
//		scanner := bufio.NewScanner(file)
//
//		for scanner.Scan() {
//			i++
//			//line := scanner.Bytes()
//			userItem := &user.JsonUser{}
//			err := easyjson.Unmarshal(scanner.Bytes(), userItem)
//			if err != nil {
//				continue
//			}
//
//			isIE := false
//			isAndroid := false
//
//			for _, browser := range userItem.Browsers {
//				anyTargetBrowser := false
//				if strings.Contains(browser, "Android") {
//					isAndroid = true
//					anyTargetBrowser = true
//				}
//
//				if strings.Contains(browser, "MSIE") {
//					isIE = true
//					anyTargetBrowser = true
//				}
//
//				if anyTargetBrowser {
//					browsers[browser] = struct{}{}
//				}
//			}
//
//			if !(isIE && isAndroid) {
//				continue
//			}
//
//			email := r.ReplaceAllString(userItem.Email, " [at] ")
//			_, _ = fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, userItem.Name, email))
//		}
//
//		_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
//	}
//var threadsCount = runtime.NumCPU()
//
//func FastSearch(out io.Writer) {
//
//	_, _ = fmt.Fprintln(out, "found users:")
//
//	browsers := map[string]struct{}{}
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	i := -1
//	scanner := bufio.NewScanner(file)
//	jsonParsed := make([]*user.JsonUser, threadsCount)
//	wg := sync.WaitGroup{}
//
//	for scanner.Scan() {
//		i++
//
//		wg.Add(1)
//		lineBytes := scanner.Bytes()
//		copiedBytes := make([]byte, len(lineBytes))
//		_ = copy(copiedBytes, lineBytes)
//		go func(data []byte, i int) {
//			defer wg.Done()
//			userItem := &user.JsonUser{}
//			err := easyjson.Unmarshal(data, userItem)
//			if err == nil {
//				jsonParsed[i] = userItem
//			} else {
//				jsonParsed[i] = nil
//			}
//		}(copiedBytes, i%threadsCount)
//
//		if (i+1)%threadsCount == 0 {
//			wg.Wait()
//			FastSearchJob(out, jsonParsed, browsers, i+1-threadsCount)
//		}
//	}
//
//	if i > -1 && (i+1)%threadsCount != 0 {
//		wg.Wait()
//		FastSearchJob(out, jsonParsed, browsers, i-i%threadsCount)
//	}
//
//	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
//}
//
//func FastSearchJob(out io.Writer, items []*user.JsonUser, browsers map[string]struct{}, startIndex int) {
//
//	for i, userItem := range items {
//		if userItem == nil {
//			continue
//		}
//		isIE := false
//		isAndroid := false
//
//		for _, browser := range userItem.Browsers {
//			anyTargetBrowser := false
//			if strings.Contains(browser, "Android") {
//				isAndroid = true
//				anyTargetBrowser = true
//			}
//
//			if strings.Contains(browser, "MSIE") {
//				isIE = true
//				anyTargetBrowser = true
//			}
//
//			if anyTargetBrowser {
//				browsers[browser] = struct{}{}
//			}
//		}
//
//		if !(isIE && isAndroid) {
//			continue
//		}
//
//		email := r.ReplaceAllString(userItem.Email, " [at] ")
//		_, _ = fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", startIndex+i, userItem.Name, email))
//	}
//}
//
//type InChData struct {
//	index int
//	data  []byte
//}
//
//type OutChData struct {
//	index int
//	data  string
//}
//
//var threadsCount = runtime.NumCPU()
//
//func FastSearch(out io.Writer) {
//
//	_, _ = fmt.Fprintln(out, "found users:")
//
//	browsers := map[string]struct{}{}
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	i := -1
//	scanner := bufio.NewScanner(file)
//
//	inCh := make(chan InChData, threadsCount)
//	outCh := make(chan OutChData, threadsCount)
//
//	var count int32 = 0
//	var endingCount int32 = 0
//
//	mu := sync.Mutex{}
//
//	for i := 0; i < threadsCount; i++ {
//		go func() {
//			for in := range inCh {
//				userItem := &user.JsonUser{}
//				err := easyjson.Unmarshal(in.data, userItem)
//				if err != nil {
//					if atomic.AddInt32(&count, 1) == int32(threadsCount) {
//						close(outCh)
//					}
//					continue
//				}
//				isIE := false
//				isAndroid := false
//
//				for _, browser := range userItem.Browsers {
//					anyTargetBrowser := false
//					if strings.Contains(browser, "Android") {
//						isAndroid = true
//						anyTargetBrowser = true
//					}
//
//					if strings.Contains(browser, "MSIE") {
//						isIE = true
//						anyTargetBrowser = true
//					}
//
//					if anyTargetBrowser {
//						mu.Lock()
//						browsers[browser] = struct{}{}
//						mu.Unlock()
//					}
//				}
//
//				if !(isIE && isAndroid) {
//					if atomic.AddInt32(&count, 1) == int32(threadsCount) {
//						close(outCh)
//					}
//					continue
//				}
//
//				email := r.ReplaceAllString(userItem.Email, " [at] ")
//				outCh <- OutChData{in.index % threadsCount, fmt.Sprintf("[%d] %s <%s>", in.index, userItem.Name, email)}
//
//				if atomic.AddInt32(&count, 1) == int32(threadsCount) {
//					close(outCh)
//				}
//
//			}
//			if atomic.AddInt32(&endingCount, 1) == int32(threadsCount) {
//				close(outCh)
//			}
//		}()
//	}
//
//	for scanner.Scan() {
//		i++
//
//		lineBytes := scanner.Bytes()
//		copiedBytes := make([]byte, len(lineBytes))
//		_ = copy(copiedBytes, lineBytes)
//		inCh <- InChData{i, copiedBytes}
//
//		if (i+1)%threadsCount == 0 {
//			result := make(map[int]string, threadsCount)
//			for out := range outCh {
//				result[out.index] = out.data
//			}
//			keys := make([]int, 0, len(result))
//
//			for k := range result {
//				keys = append(keys, k)
//			}
//			sort.Ints(keys)
//
//			for _, k := range keys {
//				_, _ = fmt.Fprintln(out, result[k])
//			}
//
//			outCh = make(chan OutChData, threadsCount)
//			count = 0
//		}
//	}
//
//	close(inCh)
//
//	if i > -1 && (i+1)%threadsCount != 0 {
//		result := make(map[int]string, threadsCount)
//		for out := range outCh {
//			result[out.index] = out.data
//		}
//		keys := make([]int, 0, len(result))
//
//		for k := range result {
//			keys = append(keys, k)
//		}
//		sort.Ints(keys)
//
//		for _, k := range keys {
//			_, _ = fmt.Fprintln(out, result[k])
//		}
//	}
//
//	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
//}

//type InChData struct {
//	index int
//	data  []byte
//}
//
//type OutChData struct {
//	index int
//	data  user.JsonUser
//}
//
//var threadsCount = runtime.NumCPU()
//
//func FastSearch(out io.Writer) {
//
//	_, _ = fmt.Fprintln(out, "found users:")
//
//	browsers := map[string]struct{}{}
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	i := -1
//	scanner := bufio.NewScanner(file)
//
//	inCh := make(chan InChData, threadsCount)
//	outCh := make(chan OutChData, threadsCount)
//
//	var count int32 = 0
//	var endingCount int32 = 0
//
//	//mu := sync.Mutex{}
//
//	for i := 0; i < threadsCount; i++ {
//		go func() {
//			for in := range inCh {
//				userItem := user.JsonUser{}
//				err := easyjson.Unmarshal(in.data, &userItem)
//				if err != nil {
//					if atomic.AddInt32(&count, 1) == int32(threadsCount) {
//						close(outCh)
//					}
//					continue
//				}
//
//				outCh <- OutChData{in.index, userItem}
//
//				//isIE := false
//				//isAndroid := false
//				//
//				//for _, browser := range userItem.Browsers {
//				//	anyTargetBrowser := false
//				//	if strings.Contains(browser, "Android") {
//				//		isAndroid = true
//				//		anyTargetBrowser = true
//				//	}
//				//
//				//	if strings.Contains(browser, "MSIE") {
//				//		isIE = true
//				//		anyTargetBrowser = true
//				//	}
//				//
//				//	if anyTargetBrowser {
//				//		mu.Lock()
//				//		browsers[browser] = struct{}{}
//				//		mu.Unlock()
//				//	}
//				//}
//				//
//				//if !(isIE && isAndroid) {
//				//	if atomic.AddInt32(&count, 1) == int32(threadsCount) {
//				//		close(outCh)
//				//	}
//				//	continue
//				//}
//				//
//				//email := r.ReplaceAllString(userItem.Email, " [at] ")
//				//outCh <- OutChData{in.index % threadsCount, fmt.Sprintf("[%d] %s <%s>", in.index, userItem.Name, email)}
//
//				if atomic.AddInt32(&count, 1) == int32(threadsCount) {
//					close(outCh)
//				}
//
//			}
//			if atomic.AddInt32(&endingCount, 1) == int32(threadsCount) {
//				close(outCh)
//			}
//		}()
//	}
//
//	for scanner.Scan() {
//		i++
//
//		lineBytes := scanner.Bytes()
//		copiedBytes := make([]byte, len(lineBytes))
//		_ = copy(copiedBytes, lineBytes)
//		inCh <- InChData{i, copiedBytes}
//
//		if (i+1)%threadsCount == 0 {
//			result := make(map[int]user.JsonUser, threadsCount)
//			for out := range outCh {
//				result[out.index] = out.data
//			}
//			keys := make([]int, 0, len(result))
//
//			for k := range result {
//				keys = append(keys, k)
//			}
//			sort.Ints(keys)
//
//			for _, k := range keys {
//				FastSearchJob(out, result[k], browsers, k)
//			}
//
//			outCh = make(chan OutChData, threadsCount)
//			count = 0
//		}
//	}
//
//	close(inCh)
//
//	if i > -1 && (i+1)%threadsCount != 0 {
//		result := make(map[int]user.JsonUser, threadsCount)
//		for out := range outCh {
//			result[out.index] = out.data
//		}
//		keys := make([]int, 0, len(result))
//
//		for k := range result {
//			keys = append(keys, k)
//		}
//		sort.Ints(keys)
//
//		for _, k := range keys {
//			FastSearchJob(out, result[k], browsers, k)
//		}
//	}
//
//	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
//}
//
//func FastSearchJob(out io.Writer, userItem user.JsonUser, browsers map[string]struct{}, index int) {
//
//	isIE := false
//	isAndroid := false
//
//	for _, browser := range userItem.Browsers {
//		anyTargetBrowser := false
//		if strings.Contains(browser, "Android") {
//			isAndroid = true
//			anyTargetBrowser = true
//		}
//
//		if strings.Contains(browser, "MSIE") {
//			isIE = true
//			anyTargetBrowser = true
//		}
//
//		if anyTargetBrowser {
//			browsers[browser] = struct{}{}
//		}
//	}
//
//	if !(isIE && isAndroid) {
//		return
//	}
//
//	email := r.ReplaceAllString(userItem.Email, " [at] ")
//	_, _ = fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", index, userItem.Name, email))
//
//}

//type InChData struct {
//	index int
//	data  []byte
//}
//
//type OutChData struct {
//	index int
//	data  user.JsonUser
//}
//
//var threadsCount = runtime.NumCPU()
//
//func FastSearch(out io.Writer) {
//
//	_, _ = fmt.Fprintln(out, "found users:")
//
//	browsers := map[string]struct{}{}
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	i := -1
//	scanner := bufio.NewScanner(file)
//
//	inCh := make(chan InChData)
//	outCh := make(chan OutChData)
//
//	var endingCount int32 = 0
//	//mu := &sync.Mutex{}
//
//	for i := 0; i < threadsCount; i++ {
//		go func() {
//			for in := range inCh {
//				userItem := user.JsonUser{}
//				err := easyjson.Unmarshal(in.data, &userItem)
//				if err != nil {
//					//continue
//				}
//
//				outCh <- OutChData{in.index, userItem}
//			}
//			if atomic.AddInt32(&endingCount, 1) == int32(threadsCount) {
//				close(outCh)
//			}
//		}()
//	}
//
//	go func() {
//		for scanner.Scan() {
//			i++
//
//			lineBytes := scanner.Bytes()
//			copiedBytes := make([]byte, len(lineBytes))
//			_ = copy(copiedBytes, lineBytes)
//			inCh <- InChData{i, copiedBytes}
//		}
//		close(inCh)
//	}()
//
//	result := make([]OutChData, 0)
//	minTarget := 0
//	maxValue := 0
//	for response := range outCh {
//		result = append(result, response)
//		if response.index > maxValue {
//			maxValue = response.index
//		}
//		if maxValue-minTarget == len(result)-1 && len(result) > 100 {
//
//			sort.Slice(result, func(i, j int) bool {
//				return result[i].index < result[j].index
//			})
//
//			for k := range result {
//
//				FastSearchJob(out, result[k].data, browsers, result[k].index)
//			}
//
//			minTarget = maxValue + 1
//			result = make([]OutChData, 0)
//		}
//	}
//
//	if len(result) > 0 {
//		sort.Slice(result, func(i, j int) bool {
//			return result[i].index < result[j].index
//		})
//		for k := range result {
//			FastSearchJob(out, result[k].data, browsers, result[k].index)
//		}
//	}
//
//	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
//}
//
//func FastSearchJob(out io.Writer, userItem user.JsonUser, browsers map[string]struct{}, index int) {
//
//	isIE := false
//	isAndroid := false
//
//	for _, browser := range userItem.Browsers {
//		anyTargetBrowser := false
//		if strings.Contains(browser, "Android") {
//			isAndroid = true
//			anyTargetBrowser = true
//		}
//
//		if strings.Contains(browser, "MSIE") {
//			isIE = true
//			anyTargetBrowser = true
//		}
//
//		if anyTargetBrowser {
//			browsers[browser] = struct{}{}
//		}
//	}
//
//	if !(isIE && isAndroid) {
//		return
//	}
//
//	email := r.ReplaceAllString(userItem.Email, " [at] ")
//	_, _ = fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", index, userItem.Name, email))
//
//}
