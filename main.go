package main

import (
	"marb.ec/corvi-backend/views"
	"marb.ec/maf/router"
	"net/http"
)

func main() {

	r := router.NewTreeRouter()

	// Category Routes
	r.Add(router.GET, "/api/categories", &views.CategoriesView{})
	r.Add(router.GET, "/api/category/:id", &views.CategoryView{})
	r.Add(router.GET, "/api/category/:id/boxes", &views.CategoryBoxesView{})

	// Boxes Routes
	r.Add(router.GET, "/api/boxes", &views.BoxesView{})
	r.Add(router.GET, "/api/box/:id", &views.BoxView{})
	r.Add(router.GET, "/api/box/:id/questions", &views.BoxQuestionsView{})
	r.Add(router.GET, "/api/box/:id/getQuestionToLearn", &views.BoxGetQuestionToLearnView{})

	// Question Routes
	r.Add(router.GET, "/api/questions", &views.QuestionsView{})
	r.Add(router.GET, "/api/question/:id", &views.QuestionView{})
	r.Add(router.PUT, "/api/question/:id/giveCorrectAnswer", &views.QuestionGiveCorrectAnswerView{})
	r.Add(router.PUT, "/api/question/:id/giveWrongAnswer", &views.QuestionGiveWrongAnswerView{})

	// Statistics Routes

	// Discovery / Cloud Routes

	// Settings Routes

	// TODO(mjb): Add Middleware
	// TODO(mjb): Restrict access to electron (via header field?)

	// Only bind to localhost for electron
	http.ListenAndServe("127.0.0.1:8080", r)

}
