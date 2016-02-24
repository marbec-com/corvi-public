package main

import (
	"log"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/corvi-backend/views"
	"marb.ec/maf/events"
	"marb.ec/maf/interfaces"
	"marb.ec/maf/router"
	"marb.ec/maf/wsnotify"
	"net/http"
	"time"
)

func main() {

	r := router.NewTreeRouter()

	go func() {
		eh := events.Events()
		for _ = range time.Tick(10 * time.Second) {
			err := eh.Publish(interfaces.Topic("categories"), nil)
			log.Printf("Publish: %v", err)
		}
	}()

	topics := []interfaces.Topic{"categories", "boxes", "questions"}
	ns := wsnotify.NewWSNotificationService(topics)

	// WS Notify Routes
	r.Add(router.GET, "/sock", ns)

	// Static Routes
	r.Add(router.GET, "/app", &controllers.IndexController{})
	r.Add(router.GET, "/app/*path", &controllers.FileController{})

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
