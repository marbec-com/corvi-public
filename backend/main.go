package main

import (
	"log"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/corvi-backend/middleware"
	"marb.ec/corvi-backend/views"
	"marb.ec/maf/requests"
	"marb.ec/maf/router"
	"marb.ec/maf/wsnotify"
)

func main() {

	// TODO(mjb): Timer at change of day to refill and refresh QuestionHeaps of all boxes

	db, err := controllers.NewDBController("data.db")
	if err != nil {
		log.Fatal(err)
	}
	controllers.InitControllerSingletons(db)

	r := router.NewTreeRouter()
	defineRoutes(r)

	// TODO(mjb): Restrict access to electron (via header field?)
	webserver := requests.NewRequestHandler(r)
	webserver.SetNotFoundHandler(&middleware.NotFoundHandler{})
	webserver.AppendGlobalPreHandler(&middleware.LogHandler{})
	webserver.PrependGlobalPreHandler(&middleware.PanicRecoveryHandler{})

	// Only bind to localhost for electron
	log.Fatal(webserver.ListenAndServe("127.0.0.1:8080"))

}

func defineRoutes(r *router.TreeRouter) {

	// WebSocket Notification Service
	ns := wsnotify.NewWSNotificationService()
	r.Add(router.GET, "/sock", ns)

	// Static Routes
	r.Add(router.GET, "/app", &controllers.IndexController{})
	r.Add(router.GET, "/app/*path", &controllers.FileController{})

	// Category Routes
	r.Add(router.GET, "/api/categories", &views.CategoriesView{})
	r.Add(router.POST, "/api/categories", &views.CategoryAddView{})
	r.Add(router.GET, "/api/category/:id", &views.CategoryView{})
	r.Add(router.PUT, "/api/category/:id", &views.CategoryUpdateView{})
	r.Add(router.DELETE, "/api/category/:id", &views.CategoryDeleteView{})
	r.Add(router.GET, "/api/category/:id/boxes", &views.CategoryBoxesView{})

	// Boxes Routes
	r.Add(router.GET, "/api/boxes", &views.BoxesView{})
	r.Add(router.POST, "/api/boxes", &views.BoxAddView{})
	r.Add(router.GET, "/api/box/:id", &views.BoxView{})
	r.Add(router.PUT, "/api/box/:id", &views.BoxUpdateView{})
	r.Add(router.DELETE, "/api/box/:id", &views.BoxDeleteView{})
	r.Add(router.GET, "/api/box/:id/questions", &views.BoxQuestionsView{})
	r.Add(router.GET, "/api/box/:id/getQuestionToLearn", &views.BoxGetQuestionToLearnView{})

	// Question Routes
	r.Add(router.GET, "/api/questions", &views.QuestionsView{})
	r.Add(router.POST, "/api/questions", &views.QuestionAddView{})
	r.Add(router.GET, "/api/question/:id", &views.QuestionView{})
	r.Add(router.PUT, "/api/question/:id", &views.QuestionUpdateView{})
	r.Add(router.DELETE, "/api/question/:id", &views.QuestionDeleteView{})
	r.Add(router.PUT, "/api/question/:id/giveCorrectAnswer", &views.QuestionGiveCorrectAnswerView{})
	r.Add(router.PUT, "/api/question/:id/giveWrongAnswer", &views.QuestionGiveWrongAnswerView{})

	// Statistics Routes
	r.Add(router.GET, "/api/stats", &views.StatsView{})

	// Discovery / Cloud Routes

	// Settings Routes
	r.Add(router.GET, "/api/settings", &views.SettingsView{})
	r.Add(router.PUT, "/api/settings", &views.SettingsUpdateView{})

}
