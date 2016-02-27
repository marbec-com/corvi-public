package main

import (
	//"log"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/corvi-backend/views"
	//"marb.ec/maf/events"
	"marb.ec/maf/router"
	"marb.ec/maf/wsnotify"
	"net/http"
	//"time"
)

func main() {

	// TODO(mjb): Singletons thread safe? especially settings!
	// TODO(mjb): Timer at change of day to refill and refresh QuestionHeaps of all boxes

	r := router.NewTreeRouter()

	/* go func() {
		eh := events.Events()
		i := 0
		for _ = range time.Tick(10 * time.Second) {
			i++
			if i%3 == 0 {
				eh.Publish(events.Topic("boxes"), nil)
				log.Println("Publish All")
			} else if i%3 == 1 {
				eh.Publish(events.Topic("box-1"), nil)
				log.Println("Publish Single")
			} else {
				eh.Publish(events.Topic("boxcat-1"), nil)
				log.Println("Publish Boxcat")
			}
		}
	}() */

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
	r.Add(router.GET, "/api/category/:id/boxes", &views.CategoryBoxesView{})

	// Boxes Routes
	r.Add(router.GET, "/api/boxes", &views.BoxesView{})
	r.Add(router.POST, "/api/boxes", &views.BoxAddView{})
	r.Add(router.GET, "/api/box/:id", &views.BoxView{})
	r.Add(router.PUT, "/api/box/:id", &views.BoxUpdateView{})
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
	r.Add(router.GET, "/api/settings", &views.SettingsView{})
	r.Add(router.PUT, "/api/settings", &views.SettingsUpdateView{})

	// TODO(mjb): Add Middleware
	// TODO(mjb): Restrict access to electron (via header field?)

	// Only bind to localhost for electron
	http.ListenAndServe("127.0.0.1:8080", r)

}
