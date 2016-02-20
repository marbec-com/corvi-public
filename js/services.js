var corviServices = angular.module('corviServices', []);

corviServices.factory('Categories', function($http, $log) {
	var CategoryService = {};
	
	CategoryService.getAll = function(callback) {
		$http.get("/api/categories/").then(function(res) {
			callback(res.data);
		}, function(res) {
			$log.error(res);
		});
	};
	
	CategoryService.getAllWithBoxes = function(callback) {
		$http.get("/api/categories/").then(function(res) {
			var categories = res.data;
			Object.keys(categories).forEach(function(key) {
				CategoryService.getBoxes(categories[key].ID, function(data) {
					categories[key].Boxes = data;
				});
			});
			callback(categories);
		}, function(res) {
			$log.error(res);
		});
	};
	
	CategoryService.get = function(catID, callback) {
		$http.get("/api/category/"+catID+"/").then(function(res) {
			callback(res.data);
		}, function(res) {
			$log.error(res);
		});
	};
	
	CategoryService.getBoxes = function(catID, callback) {
		$http.get("/api/category/"+catID+"/boxes/").then(function(res) {
			callback(res.data);
		}, function(res) {
			$log.error(res);
		});
	};
	
	CategoryService.update = function(category, callback) {
		delete(category.Boxes);
		$log.debug(category);
	};
	
	return CategoryService;
});

corviServices.factory('Boxes', function($http, $log) {
	var BoxService = {};
	
	BoxService.get = function(boxID, success, error) {
		$http.get("/api/box/"+boxID+"/").then(function(res) {
			success(res.data);
		}, function(res) {
			$log.error(res);
			error(res);
		});
	};
	
	return BoxService;

});

corviServices.factory('Questions', function($http, $log) {
	var QuestionService = {};
	
	QuestionService.getQuestionToLearn = function(boxID, question, finished, error) {
		$http.get("/api/box/"+boxID+"/getQuestionToLearn").then(function(res) {
			if(res.status == 200) {
				question(res.data);
			}else{ // No more questions
				finished();
			}
		}, function(res) {
			$log.error(res);
			error(res);
		});
	};
	
	QuestionService.giveCorrectAnswer = function(questionID, success, error) {
		$http.put("/api/question/"+questionID+"/giveCorrectAnswer").then(function(res) {
			success();
		}, function(res) {
			$log.error(res);
			error();
		});
	};
	
	QuestionService.giveWrongAnswer = function(questionID, success, error) {
		$http.put("/api/question/"+questionID+"/giveWrongAnswer").then(function(res) {
			success();
		}, function(res) {
			$log.error(res);
			error();
		});
	};
	
	return QuestionService;
});