package main

type TaskData struct {
	Data struct {
		ProblemsetQuestionList struct {
			Total     int        `json:"total"`
			Questions []Question `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

type Question struct {
	AcRate     float64 `json:"acRate"`
	Difficulty string  `json:"difficulty"`
	PaidOnly   bool    `json:"paidOnly"`
	Title      string  `json:"title"`
	TitleSlug  string  `json:"titleSlug"`
	TopicTags  []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"topicTags"`
}

type NumData struct {
	Data struct {
		ProblemsetQuestionList struct {
			Total int `json:"total"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

type DailyData struct {
	ActiveDailyCodingChallengeQuestion DailyTask `json:"activeDailyCodingChallengeQuestion"`
}

type DailyTask struct {
	Date       string        `json:"date"`
	UserStatus string        `json:"userStatus"`
	Link       string        `json:"link"`
	Question   DailyQuestion `json:"question"`
}

type DailyQuestion struct {
	AcRate             float64 `json:"acRate"`
	Difficulty         string  `json:"difficulty"`
	FrontendQuestionId string  `json:"frontendQuestionId"`
	IsFavor            bool    `json:"isFavor"`
	PaidOnly           bool    `json:"paidOnly"`
	Title              string  `json:"title"`
	TitleSlug          string  `json:"titleSlug"`
	HasVideoSolution   bool    `json:"hasVideoSolution"`
	HasSolution        bool    `json:"hasSolution"`
	TopicTags          []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
		Slug string `json:"slug"`
	} `json:"topicTags"`
}

type UserData struct {
	MatchedUser UserStats `json:"matchedUser"`
}

type UserStats struct {
	Username    string `json:"username"`
	SubmitStats struct {
		AcSubmissionNum []struct {
			Difficulty  string `json:"difficulty"`
			Count       int    `json:"count"`
			Submissions int    `json:"submissions"`
		} `json:"acSubmissionNum"`
	} `json:"submitStats"`
}

type Payload struct {
	Query         string    `json:"query"`
	Variables     Variables `json:"variables"`
	OperationName string    `json:"operationName"`
}
type Filters struct {
}
type Variables struct {
	CategorySlug string  `json:"categorySlug"`
	Skip         int     `json:"skip"`
	Limit        int     `json:"limit"`
	Filters      Filters `json:"filters"`
}
