/* global angular */
var corviServices = angular.module('corviServices', []);

corviServices.factory('Categories', function($http, $log) {
	var CategoryService = {};
	
	CategoryService.CategoriesAll = [];
	CategoryService.CategoriesByID = {};
	
	CategoryService.Update = function() {
		$http.get("/api/categories/").then(function(res) {
			// Clear array and object, but preserve reference
			CategoryService.CategoriesAll.length = 0;
			ClearObject(CategoryService.CategoriesByID);
			
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
		$http.get("/api/category/"+catID+"/").then(function(res) {
			// Update individual entry in CategoryService.CategoriesByID, preserve reference
			// This also updates the same object in CategoryService.CategoriesAll
			var newCategory = angular.copy(res.data);
			angular.copy(newCategory, CategoryService.CategoriesByID[catID]);
		}, function(res) {
			$log.error(res);
		});	
	};
	
	return CategoryService;
});

corviServices.factory('Boxes', function($http, $log) {
	var BoxService = {};
	
	BoxService.BoxesByID = {};
	BoxService.BoxesByCatID = {};
	BoxService.BoxesAll = [];
	
	BoxService.Update = function() {
		$http.get("/api/boxes/").then(function(res) {
			// Clear array and objects, but preserve reference
			BoxService.BoxesAll.length = 0;
			ClearObject(BoxService.BoxesByID);
			ClearObject(BoxService.BoxesByCatID);
			
			// Fill with new data
			for (var i = 0; i < res.data.length; i++) {
				var newBox = angular.copy(res.data[i]);
				BoxService.BoxesAll.push(newBox);
				BoxService.BoxesByID[newBox.ID] = newBox;
				if(!(newBox.Category.ID in BoxService.BoxesByCatID)) {
					BoxService.BoxesByCatID[newBox.Category.ID] = [];
				}
				BoxService.BoxesByCatID[newBox.Category.ID].push(newBox);			
			}
		}, function(res) {
			$log.error(res);
		}); 
	};
	
	BoxService.UpdateSingle = function(boxID) {
		$http.get("/api/box/"+boxID+"/").then(function(res) {
			var newBox = angular.copy(res.data);
			angular.copy(newBox, BoxService.BoxesByID[boxID]);
		}, function(res) {
			$log.error(res);
		});	
	};
	
	BoxService.UpdateCategory = function(catID) {
		$http.get("/api/category/"+catID+"/boxes/").then(function(res) {
			// Clear array of BoxesByCatID object or create a new one, but preserve reference
			if (!(catID in BoxService.BoxesByCatID)) {
				BoxService.BoxesByCatID[catID] = [];
			}else{
				BoxService.BoxesByCatID[catID].length = 0
			}
			
			// Fill with new data
			for (var i = 0; i < res.data.length; i++) {
				var newBox = angular.copy(res.data[i]);		
				// Update by id
				angular.copy(newBox, BoxService.BoxesByID[newBox.ID]);
				// Add to bycatid
				BoxService.BoxesByCatID[catID].push(newBox);
			}
		}, function(res) {
			$log.error(res);
		}); 	
	};
	
	return BoxService;

});

corviServices.factory('Questions', function($http, $log) {
	var QuestionService = {};
	
	QuestionService.Update = function() {
		
	};
	
	QuestionService.UpdateSingle = function(qID) {
		
	};
	
	QuestionService.UpdateBox = function(boxID) {
		
	};
	
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
	
	var parseID = function(message) {
		var parts = message.split("-");
		if(parts.length != 2) {
			return -1
		}
		var id = parseInt(parts[1], 10);
		if (isNaN(id)) {
			return -1
		}
		return id
	};
	
	NotifyService.onOpen = function(res) {
		$log.debug("Websocket connection opened: ", res);
	};
	
	NotifyService.onMessage = function(res) {
		$log.debug("Incoming Websocket notification: ", res);
		if(res.data == "categories") {
			Categories.Update();
		}else if(StartsWith(res.data, "category-")) {
			var catID = parseID(res.data);
			if(catID === -1) {
				$log.error("Invalid ID: ", catID);
				return
			}
			Categories.UpdateSingle(catID);
		}else if(res.data == "boxes") {
			Boxes.Update();
		}else if(StartsWith(res.data, "box-")) {
			var boxID = parseID(res.data);
			if(boxID === -1) {
				$log.error("Invalid ID: ", boxID);
				return
			}
			Boxes.UpdateSingle(boxID);
		}else if(StartsWith(res.data, "boxcat-")) {
			var catID = parseID(res.data);
			if(catID === -1) {
				$log.error("Invalid ID: ", catID);
				return
			}
			Boxes.UpdateCategory(catID);
		}else if(res.data == "questions") {
			Questions.Update();
		}else if(StartsWith(res.data, "question-")) {
			var qID = parseID(res.data);
			if(qID === -1) {
				$log.error("Invalid ID: ", qID);
				return
			}
			Questions.UpdateSingle(qID);
		}else if(StartsWith(res.data, "questionbox-")) {
			var boxID = parseID(res.data);
			if(boxID === -1) {
				$log.error("Invalid ID: ", boxID);
				return
			}
			Questions.UpdateBox(boxID);
		}else{
			$log.debug("Unknown Websocket notification: ", res);
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