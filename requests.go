package main

var greet = `Hello, I am a utility LeetCode bot.

You can:
‣ _view_  your LeetCode stats,
‣ _solve_  the daily task,
‣ _pick_  a random one,
‣ _read_  about the bot.
Your choice?`

var aboutText = `This is a simple Telegram bot implemented in _Go_ for an internship in VK.
It levers the _Telegram Bot API_.
Are you interested in the _source code_?`

var requestTask = `query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
 problemsetQuestionList: questionList(
	categorySlug: $categorySlug
	limit: $limit
	skip: $skip
	filters: $filters
	) {
		total: totalNum
		questions: data {
			acRate
			difficulty
			paidOnly: isPaidOnly
			title
			titleSlug
			topicTags {
				name
				slug
			}
		}
	}
}`

var requestNum = `query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
 problemsetQuestionList: questionList(
	categorySlug: $categorySlug
	limit: $limit
	skip: $skip
	filters: $filters
	) {
		total: totalNum
	}
}`

var requestDaily = `query questionOfToday {
  activeDailyCodingChallengeQuestion {
      date
      userStatus
      link
      question {
          acRate
          difficulty
          frontendQuestionId: questionFrontendId
          isFavor
          paidOnly: isPaidOnly
          title
          titleSlug
          hasVideoSolution
          hasSolution
          topicTags {
              name
              id
              slug
          }
      }
  }
}`

var requestUser = `query getUserProfile($username: String!) {
	matchedUser(username: $username) {
		username
		submitStats: submitStatsGlobal {
			acSubmissionNum {
				difficulty
				count
				submissions
			}
		}
	}
}`
