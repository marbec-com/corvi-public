var corviServices = angular.module('corviServices', []);

corviServices.factory('Categories', function($http, $log, Boxes) {
	var CategoryService = {};
	
	CategoryService.Categories = [];
	CategoryService.CategoriesRaw = {};
	
	var RawToArray = function() {
		CategoryService.Categories.length = 0;
		Object.keys(CategoryService.CategoriesRaw).forEach(function(key) {
			CategoryService.Categories.push(CategoryService.CategoriesRaw[key]);
		});
	};
	
	var AssignBoxes = function() {
		for (var i = 0; i < Boxes.Boxes.length; i++) {
			CategoryService.CategoriesRaw[Boxes.Boxes[i].Category.ID].Boxes.push(Boxes.Boxes[i]);
		}
	};
	
	CategoryService.Update = function(callback) {
		$http.get("/api/categories/").then(function(res) {
			CategoryService.CategoriesRaw = {};
			for (var i = 0; i < res.data.length; i++) {
				CategoryService.CategoriesRaw[res.data[i].ID] = res.data[i];
				CategoryService.CategoriesRaw[res.data[i].ID].Boxes = [];
			}
			AssignBoxes();
			RawToArray();
			$log.debug(CategoryService.Categories);
			if (callback) {
				callback();
			}
		}, function(res) {
			$log.error(res);
		});
	};
	
	CategoryService.UpdateCategory = function(catID, callback) {
		$http.get("/api/category/"+catID+"/").then(function(res) {
			CategoryService.CategoriesRaw[res.data.ID] = res.data;
			CategoryService.CategoriesRaw[res.data.ID].Boxes = [];
			AssignBoxes();
			RawToArray();
			callback();
		}, function(res) {
			$log.error(res);
		});
	};
	
	return CategoryService;
});

// TODO save boxes in arrays by catID, same for questions
/*
boxes = {}
boxes[catID] = [a, b, c]...
*/
corviServices.factory('Boxes', function($http, $log) {
	var BoxService = {};
	
	BoxService.Boxes = [];
	BoxService.BoxesRaw = {};
	
	var RawToArray = function() {
		BoxService.Boxes.length = 0;
		Object.keys(BoxService.BoxesRaw).forEach(function(key) {
			BoxService.Boxes.push(BoxService.BoxesRaw[key]);
		});
	};
	
	BoxService.Update = function(callback) {
		$http.get("/api/boxes/").then(function(res) {
			BoxService.Boxes = res.data;
			BoxService.BoxesRaw = {};
			for (var i = 0; i < res.data.length; i++) {
				BoxService.BoxesRaw[res.data[i].ID] = res.data[i];
				// TODO Update Questions
			}
			RawToArray();
			$log.debug(BoxService.Boxes);
			$log.debug(BoxService.BoxesRaw);
			callback();
		}, function(res) {
			$log.error(res);
		});
	};
	
	BoxService.UpdateBox = function(boxID, callback) {
		$http.get("/api/box/"+boxID+"/").then(function(res) {
			BoxService.BoxesRaw[res.data.ID] = res.data;
			// TODO Update Questions
			RawToArray();
			callback();
		}, function(res) {
			$log.error(res);
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

corviServices.factory('Notify', function($http, $log, Categories, Boxes, Questions) {
	var NotifyService = {};
	
	$log.debug("NotifyService");
	
	NotifyService.onOpen = function(res) {
		$log.debug("Open: ", res);
	};
	
	NotifyService.onMessage = function(res) {
		switch(res.data) {
			case "categories":
				$log.debug("Update categories");
				Categories.Update();
				// load categories
				break;
			case "boxes":
				$log.debug("Update boxes");
				// load boxes
				break;
			case "questions":
				$log.debug("Update questions");
				// load questions
				break;
			default:
				$log.debug("Unknown message: ", res);
		}
	};
	
	NotifyService.onError = function(res) {
		$log.debug("Error: ", res);
	};
	
	NotifyService.onClose = function(res) {
		$log.debug("Close: ", res);
	};
	
	NotifyService.connect = function() {
		try {
			NotifyService.sock = new WebSocket("ws://127.0.0.1:8080/sock");
			$log.debug("Websocket - status: " + NotifyService.sock.readyState);
			NotifyService.sock.onopen = NotifyService.onOpen;
			NotifyService.sock.onmessage = NotifyService.onMessage;
			NotifyService.sock.onerror = NotifyService.onError;
			NotifyService.sock.onclose = NotifyService.onClose;
		} catch(exception) {
			$log.error(exception);
		}
	};
	
	NotifyService.Destroy = function() {
		$log.debug("Close ws");
		NotifyService.onClose = function() {};
		NotifyService.sock.close();	
	};
		
	return NotifyService;
});