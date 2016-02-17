package main

import (
	"marb.ec/corvi-backend/controllers"
	"marb.ec/maf/router"
	"net/http"
)

func main() {

	r := router.NewTreeRouter()

	// Category Routes
	r.Add(router.GET, "/api/categories", &controllers.CategoryController{})
	r.Add(router.GET, "/api/category/:id", &controllers.CategoryController{})
	r.Add(router.GET, "/api/category/:id/boxes", &controllers.CategoryController{})

	// Boxes Routes
	r.Add(router.GET, "/api/boxes", &controllers.CategoryController{})
	r.Add(router.GET, "/api/box/:id", &controllers.CategoryController{})
	r.Add(router.GET, "/api/box/:id/questions", &controllers.CategoryController{})
	r.Add(router.GET, "/api/box/:id/getQuestionToLearn", &controllers.CategoryController{})

	// Question Routes
	r.Add(router.GET, "/api/questions", &controllers.CategoryController{})
	r.Add(router.GET, "/api/question/:id", &controllers.CategoryController{})
	r.Add(router.PUT, "/api/question/:id/giveCorrectAnswer", &controllers.CategoryController{})
	r.Add(router.PUT, "/api/question/:id/giveWrongAnswer", &controllers.CategoryController{})

	// TODO(mjb): Add Middleware

	http.ListenAndServe(":8080", r)

}
