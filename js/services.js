/* global angular */
var corviServices = angular.module('corviServices', []);

corviServices.factory('Categories', function($http, $log) {
	var CategoryService = {};
	
	CategoryService.CategoriesAll = [];
	CategoryService.CategoriesByID = {};
	
	var clearObject = function(obj) {
		for (var key in obj){
			if (obj.hasOwnProperty(key)){
				delete obj[key];
			}
		}
	}
	
	CategoryService.Update = function() {
		$http.get("/api/categories/").then(function(res) {
			// Clear array and object, but preserve reference
			CategoryService.CategoriesAll.length = 0;
			clearObject(CategoryService.CategoriesByID);
			
			// Fill with new data
			for (var i = 0; i < res.data.length; i++) {
				var newCategory = angular.copy(res.data[i]);
				CategoryService.CategoriesAll.push(newCategory);
				CategoryService.CategoriesByID[newCategory.ID] = newCategory;				
			}
				
			$log.debug(CategoryService.CategoriesAll, CategoryService.CategoriesByID);
		}, function(res) {
			$log.error(res);
		}); 
	};
	
	CategoryService.UpdateSingle = function(catID) {
		var id = parseInt(catID, 10);
		if (isNaN(id)) {
			$log.error("Invalid catID!");
			return
		}
	
		$http.get("/api/category/"+id+"/").then(function(res) {
			// Update individual entry in CategoryService.CategoriesByID, preserve reference
			// This also updates the same object in CategoryService.CategoriesAll
			$log.debug("Update category: ", id, res.data);
			
			var newCategory = angular.copy(res.data);
			angular.copy(newCategory, CategoryService.CategoriesByID[id]);
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
	
	/* BoxService.Boxes = [];
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
	}; */
	
	return BoxService;

});

corviServices.factory('Questions', function($http, $log) {
	var QuestionService = {};
	
	/* QuestionService.getQuestionToLearn = function(boxID, question, finished, error) {
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
	}; */
	
	return QuestionService;
});

corviServices.factory('Notify', function($http, $log, Categories, Boxes, Questions) {
	var NotifyService = {};
	
	$log.debug("NotifyService");
	
	NotifyService.onOpen = function(res) {
		$log.debug("Open: ", res);
	};
	
	NotifyService.onMessage = function(res) {
		$log.debug("Incoming message: ", res);
		switch(res.data) {
			case "categories":
				$log.debug("Update categories");
				Categories.Update();
				// load categories
				break;
			case "category-1":
				Categories.UpdateSingle(1);
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