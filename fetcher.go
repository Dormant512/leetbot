package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func HandleTask(dif string) (string, string) {
	var respData TaskData
	var numData NumData
	limit := 100

	fmt.Println(dif)

	data := Payload{
		requestNum,
		Variables{
			CategorySlug: "",
			Skip:         0,
			Limit:        1,
			Filters:      Filters{},
		},
		"problemsetQuestionList",
	}
	payloadBytes, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		log.Println(errMarshal.Error())
	}
	body := bytes.NewReader(payloadBytes)

	req, errReq := http.NewRequest("POST", "https://leetcode.com/graphql", body)
	if errReq != nil {
		log.Println(errReq.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	resp, errPost := http.DefaultClient.Do(req)
	if errPost != nil {
		log.Println(errPost.Error())
	}
	defer resp.Body.Close()

	respBody, errParse := io.ReadAll(resp.Body)
	if errParse != nil {
		log.Println(errParse.Error())
	}

	err := json.Unmarshal(respBody, &numData)
	if err != nil {
		log.Println(err.Error())
	}

	totalTasks := numData.Data.ProblemsetQuestionList.Total
	var newData TaskData
	ok := false
	for {
		ok, newData = TakeBinCheckDif(totalTasks, limit, respData, dif)
		if ok {
			break
		}
	}
	n := len(newData.Data.ProblemsetQuestionList.Questions)
	fmt.Println(dif, "questions to pick from:", n)
	taskIdx := rand.Intn(n)
	thisTask := newData.Data.ProblemsetQuestionList.Questions[taskIdx]

	taskText := "Random *" + strings.ToLower(thisTask.Difficulty) + "* task\n\n" + thisTask.Title
	taskText += "\n\n*Difficulty*: " + strings.ToLower(thisTask.Difficulty)
	taskText += "\n*Acceptance rate*: " + strconv.FormatFloat(thisTask.AcRate, 'f', 0, 64)
	taskText += "%\n*Premium only*: " + strconv.FormatBool(thisTask.PaidOnly)
	taskText += "\n" + "*Tags*:"
	for _, val := range thisTask.TopicTags {
		taskText += "\n‣ " + val.Name
	}
	thisUrl := "https://leetcode.com/problems/" + thisTask.TitleSlug + "/"
	return taskText, thisUrl
}

func TakeBinCheckDif(totalTasks, limit int, respData TaskData, dif string) (bool, TaskData) {
	totalBins := int(math.Ceil(float64(totalTasks) / float64(limit)))
	thisBin := rand.Intn(totalBins + 1)

	//fmt.Println("Bin", thisBin, "chosen")
	//fmt.Println("Len of resp", len(respData.Data.ProblemsetQuestionList.Questions))

	data := Payload{
		requestTask,
		Variables{
			CategorySlug: "",
			Skip:         thisBin * limit,
			Limit:        limit, // random in range [0, ceil(float(2673/25))+1)
			Filters:      Filters{},
		},
		"problemsetQuestionList",
	}
	payloadBytes, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		log.Println(errMarshal.Error())
	}
	body := bytes.NewReader(payloadBytes)

	req, errReq := http.NewRequest("POST", "https://leetcode.com/graphql", body)
	if errReq != nil {
		log.Println(errReq.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	resp, errPost := http.DefaultClient.Do(req)
	if errPost != nil {
		log.Println(errPost.Error())
	}
	defer resp.Body.Close()

	respBody, errParse := io.ReadAll(resp.Body)
	if errParse != nil {
		log.Println(errParse.Error())
	}

	err := json.Unmarshal(respBody, &respData)
	if err != nil {
		log.Println(err.Error())
	}

	ans := false
	respCopy := respData
	respCopy.Data.ProblemsetQuestionList.Questions = make([]Question, 0)

	for _, task := range respData.Data.ProblemsetQuestionList.Questions {
		if strings.ToLower(task.Difficulty) == dif {
			ans = true
			respCopy.Data.ProblemsetQuestionList.Questions = append(respCopy.Data.ProblemsetQuestionList.Questions, task)
		}
	}
	//if ans {
	//	fmt.Println("copy", len(respCopy.Data.ProblemsetQuestionList.Questions))
	//	fmt.Println("orig", len(respData.Data.ProblemsetQuestionList.Questions))
	//}
	return ans, respCopy
}
