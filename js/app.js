var corviApp = angular.module('corviApp', ['ngRoute']);

corviApp.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		templateUrl: 'sites/study_boxes.html',
		controller: 'studyBoxController',
		navActive: 'study'
	})
	.when('/study/:box/finished', {
		templateUrl: 'sites/study_finished.html',
		controller: 'studyFinishedController',
		navActive: 'study'
	})
	.when('/study/:box', {
		templateUrl: 'sites/study_question.html',
		controller: 'studyQuestionController',
		navActive: 'study'
	})
	.when('/manage', {
		templateUrl: 'sites/manage.html',
		controller: 'mainController',
		navActive: 'manage'
	})
	.when('/discover', {
		templateUrl: 'sites/discover.html',
		controller: 'mainController',
		navActive: 'discover'
	})
	.when('/stats', {
		templateUrl: 'sites/stats.html',
		controller: 'mainController',
		navActive: 'stats'
	})
	.when('/settings', {
		templateUrl: 'sites/settings.html',
		controller: 'mainController',
		navActive: 'settings'
	}).otherwise({
		redirectTo: '/'
	});
});

// Move $http in service

corviApp.controller('studyBoxController', function($scope, $http, $log) {
	$scope.categories = {};
	$http.get("/api/categories").then(function(res) {
		$scope.categories = res.data;
		Object.keys($scope.categories).forEach(loadBoxForCategory);
	}, function(res) {
		$log.error(res);
	});
	
	function loadBoxForCategory(catID) {
		$http.get("/api/category/" + $scope.categories[catID].ID + "/boxes").then(function(res) {
			$scope.categories[catID].Boxes = res.data;
		}, function(res) {
			$log.error(res);
		});
	}
	
});

corviApp.controller('studyFinishedController', function($scope, $routeParams, $location, $log, $http) {
	
	var id = parseInt($routeParams.box, 10)
	if (isNaN(id)) {
		$log.error("Invalid ID!");
	}
	
	$scope.box = {};
	
	$http.get("/api/box/"+id+"/").then(function(res) {
		$scope.box = res.data;
	}, function(res) {
		$log.error(res);
		$location.path("/");
	});
	
	$scope.getBack = function() {
		$location.path("/");
	};
	
});

corviApp.controller('studyQuestionController', function($scope, $routeParams, $log, $location, $http) {
	
	var boxId = parseInt($routeParams.box, 10)
	if (isNaN(boxId)) {
		$log.error("Invalid ID!");
	}

	loadNewQuestion(boxId);	
	
	$scope.showSolution = function() {
		$scope.answered = true;	
	};
	
	$scope.giveCorrectAnswer = function() {
		// Save correct answer
		$http.put("/api/question/"+$scope.question.ID+"/giveCorrectAnswer").then(function(res) {
			// Load next question
			loadNewQuestion(boxId);	
		}, function(res) {
			$log.error(res);
			$location.path("/");
		});
	};
	
	$scope.giveWrongAnswer = function() {
		// Save correct answer
		$http.put("/api/question/"+$scope.question.ID+"/giveWrongAnswer").then(function(res) {
			// Load next question
			loadNewQuestion(boxId);	
		}, function(res) {
			$log.error(res);
			$location.path("/");
		});
	};
	
	function loadNewQuestion(boxID) {
		$scope.answered = false;
		$scope.question = {};
		
		$http.get("/api/box/"+boxID+"/getQuestionToLearn").then(function(res) {
			if(res.status == 200) {
				$scope.question = res.data;
			}else{ // No more questions
				$location.path("/study/"+boxID+"/finished");
			}
		}, function(res) { // Error
			$log.error(res);
			$location.path("/");
		});
	};
});

corviApp.controller('mainController', function($scope) {
	$scope.message = 'Hello Angular!';
});

corviApp.directive('mainNavigation', function() {
	return {
		retrict: 'E',
		templateUrl: 'partials/mainNavigation.html',
		controller: 'mainNavigationController'
	}
});

corviApp.controller('mainNavigationController', function($scope, $route) {
	$scope.$route = $route;
	
	$scope.studyBadge = 3;
});