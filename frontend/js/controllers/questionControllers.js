(function() {

var corviApp = angular.module('corviQuestionControllers', ['ngRoute', 'corviServices']);

corviApp.controller('questionsController', function($scope, $routeParams, $log, $location, Questions, Boxes) {
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}
	
	$scope.compareToNow = function(time) {
		if (Date.parse(time) > Date.now()) {
			return true
		}
		return false
	}
	
	$scope.box = Boxes.BoxesByID[id];
	$scope.boxID = id;
	
	$scope.questions = Questions.QuestionsByBoxID[id];
});

corviApp.controller('questionEditController', function($scope, $routeParams, $log, $location, Questions, Boxes) {
	var id = parseInt($routeParams.question, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}
	
	$scope.question = angular.copy(Questions.QuestionsByID[id]);
	$scope.orgBoxID = $scope.question.BoxID;
	$scope.boxes = Boxes.BoxesAll;
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage/box/"+$scope.orgBoxID+"/questions");
	};
	
	$scope.save = function() {
		Questions.Update(id, $scope.question, function() {
			$location.path("/manage/box/"+$scope.orgBoxID+"/questions");
		}, function(err) {
			$scope.error = err;
		});
	};
});

corviApp.controller('questionAddController', function($scope, $routeParams, $log, $location, Questions, Boxes) {
	var boxID = parseInt($routeParams.box, 10);
	if (isNaN(boxID)) {
		$log.error("Invalid BoxID!");
		return
	}
	
	$scope.question = {};
	$scope.question.BoxID = boxID;
	$scope.boxes = Boxes.BoxesAll;
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage/box/"+boxID+"/questions");
	};
	
	$scope.save = function() {
		Questions.Add($scope.question, function(question) {
			$location.path("/manage/box/"+question.BoxID+"/questions");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

corviApp.controller('questionDeleteController', function($scope, $routeParams, $log, $location, Questions, Boxes) {
	var id = parseInt($routeParams.question, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}	
	
	$scope.question = angular.copy(Questions.QuestionsByID[id]);
	var BoxID = $scope.question.BoxID;
	$scope.box = angular.copy(Boxes.BoxesByID[BoxID]);
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage/box/"+BoxID+"/questions");
	};
	
	$scope.submit = function() {
		Questions.Delete(id, function() {
			$location.path("/manage/box/"+BoxID+"/questions");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

})();