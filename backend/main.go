package main

import (
	"github.com/facebookgo/inject"
	"log"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/corvi-backend/middleware"
	"marb.ec/corvi-backend/views"
	"marb.ec/maf/requests"
	"marb.ec/maf/router"
	"marb.ec/maf/wsnotify"
)

const (
	databaseFile string = "data.db"
	settingsFile string = "settings.yml"
)

// TODO(mjb): Timer at change of day to refill and refresh QuestionHeaps of all boxes

func main() {

	g := &inject.Graph{}

	// Initialize initializer
	Initializer, err := initInitializer(g)
	if err != nil {
		log.Fatal("Error while initializing Initializer", err)
	}

	// Initialize Router, views and routes
	r := router.NewTreeRouter()
	if err = initViewRoutes(r, g); err != nil {
		log.Fatal("Error while initializing views and routes", err)
	}

	// Initialize services
	if err = initServices(g); err != nil {
		log.Fatal("Error while initializing serivces", err)
	}

	// Initialize controllers
	if err = initControllers(g); err != nil {
		log.Fatal("Error while initializing controllers", err)
	}

	// Dependency injection
	if err := g.Populate(); err != nil {
		log.Fatal(err)
	}

	// Create tables & build heap cache
	Initializer.Init()
	Initializer.PopulateDummyData()

	// TODO(mjb): Restrict access to electron (via header field?)
	webserver := requests.NewRequestHandler(r)
	webserver.SetNotFoundHandler(&middleware.NotFoundHandler{})
	webserver.AppendGlobalPreHandler(&middleware.LogHandler{})
	webserver.PrependGlobalPreHandler(&middleware.PanicRecoveryHandler{})

	// Start Webserver
	log.Fatal(webserver.ListenAndServe("127.0.0.1:8080"))

}

func initInitializer(g *inject.Graph) (*controllers.Initializer, error) {
	i := &controllers.Initializer{}

	return i, g.Provide(
		&inject.Object{Value: i},
	)
}

func initServices(g *inject.Graph) error {

	// Init SettingsService
	settingsFileName := controllers.GenerateUserDataPath(settingsFile)
	s, err := controllers.NewYAMLSettingsService(settingsFileName)
	if err != nil {
		return err
	}

	// Init DatabaseService
	dbFile := controllers.GenerateUserDataPath(databaseFile)
	db, err := controllers.NewSQLiteDBService(dbFile)
	if err != nil {
		return err
	}

	return g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: s},
	)
}

func initControllers(g *inject.Graph) error {

	b := controllers.NewBoxController()
	c := controllers.NewCategoryController()
	q := controllers.NewQuestionController()
	st := controllers.NewStatsController()

	return g.Provide(
		&inject.Object{Value: b},
		&inject.Object{Value: c},
		&inject.Object{Value: q},
		&inject.Object{Value: st},
	)
}

func initViewRoutes(r *router.TreeRouter, g *inject.Graph) error {

	// WebSocket Notification Service
	ns := wsnotify.NewWSNotificationService()
	r.Add(router.GET, "/sock", ns)

	// Static Routes
	r.Add(router.GET, "/app", &controllers.IndexController{})
	r.Add(router.GET, "/app/*path", &controllers.FileController{})

	// Category Routes
	categoriesView := &views.CategoriesView{}
	r.Add(router.GET, "/api/categories", categoriesView)
	categoryAddView := &views.CategoryAddView{}
	r.Add(router.POST, "/api/categories", categoryAddView)
	categoryView := &views.CategoryView{}
	r.Add(router.GET, "/api/category/:id", categoryView)
	categoryUpdateView := &views.CategoryUpdateView{}
	r.Add(router.PUT, "/api/category/:id", categoryUpdateView)
	categoryDeleteView := &views.CategoryDeleteView{}
	r.Add(router.DELETE, "/api/category/:id", categoryDeleteView)
	categoryBoxesView := &views.CategoryBoxesView{}
	r.Add(router.GET, "/api/category/:id/boxes", categoryBoxesView)

	// Boxes Routes
	boxesView := &views.BoxesView{}
	r.Add(router.GET, "/api/boxes", boxesView)
	boxesAddView := &views.BoxAddView{}
	r.Add(router.POST, "/api/boxes", boxesAddView)
	boxView := &views.BoxView{}
	r.Add(router.GET, "/api/box/:id", boxView)
	boxUpdateView := &views.BoxUpdateView{}
	r.Add(router.PUT, "/api/box/:id", boxUpdateView)
	boxDeleteView := &views.BoxDeleteView{}
	r.Add(router.DELETE, "/api/box/:id", boxDeleteView)
	boxQuestionsView := &views.BoxQuestionsView{}
	r.Add(router.GET, "/api/box/:id/questions", boxQuestionsView)
	boxGetQuestionToLearnView := &views.BoxGetQuestionToLearnView{}
	r.Add(router.GET, "/api/box/:id/getQuestionToLearn", boxGetQuestionToLearnView)

	// Question Routes
	questionsView := &views.QuestionsView{}
	r.Add(router.GET, "/api/questions", questionsView)
	questionAddView := &views.QuestionAddView{}
	r.Add(router.POST, "/api/questions", questionAddView)
	questionView := &views.QuestionView{}
	r.Add(router.GET, "/api/question/:id", questionView)
	questionUpdateView := &views.QuestionUpdateView{}
	r.Add(router.PUT, "/api/question/:id", questionUpdateView)
	questionDeleteView := &views.QuestionDeleteView{}
	r.Add(router.DELETE, "/api/question/:id", questionDeleteView)
	questionGiveCorrectAnswerView := &views.QuestionGiveCorrectAnswerView{}
	r.Add(router.PUT, "/api/question/:id/giveCorrectAnswer", questionGiveCorrectAnswerView)
	questionGiveWrongAnswerView := &views.QuestionGiveWrongAnswerView{}
	r.Add(router.PUT, "/api/question/:id/giveWrongAnswer", questionGiveWrongAnswerView)

	// Statistics Routes
	statsView := &views.StatsView{}
	r.Add(router.GET, "/api/stats", statsView)

	// TODO(mjb): Discovery / Cloud Routes

	// Settings Routes
	settingsView := &views.SettingsView{}
	r.Add(router.GET, "/api/settings", settingsView)

	settingsUpdateView := &views.SettingsUpdateView{}
	r.Add(router.PUT, "/api/settings", settingsUpdateView)

	return g.Provide(
		&inject.Object{Value: categoriesView},
		&inject.Object{Value: categoryAddView},
		&inject.Object{Value: categoryView},
		&inject.Object{Value: categoryUpdateView},
		&inject.Object{Value: categoryDeleteView},
		&inject.Object{Value: categoryBoxesView},

		&inject.Object{Value: boxesView},
		&inject.Object{Value: boxesAddView},
		&inject.Object{Value: boxView},
		&inject.Object{Value: boxUpdateView},
		&inject.Object{Value: boxDeleteView},
		&inject.Object{Value: boxQuestionsView},
		&inject.Object{Value: boxGetQuestionToLearnView},

		&inject.Object{Value: questionsView},
		&inject.Object{Value: questionAddView},
		&inject.Object{Value: questionView},
		&inject.Object{Value: questionUpdateView},
		&inject.Object{Value: questionDeleteView},
		&inject.Object{Value: questionGiveCorrectAnswerView},
		&inject.Object{Value: questionGiveWrongAnswerView},

		&inject.Object{Value: statsView},

		&inject.Object{Value: settingsView},
		&inject.Object{Value: settingsUpdateView},
	)
}
