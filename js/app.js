var corviApp = angular.module('corviApp', ['ngRoute', 'corviServices']);

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

corviApp.controller('studyBoxController', function($scope, Categories) {
	$scope.categories = {};
	
	Categories.getAllWithBoxes(function(data) {
		$scope.categories = data;
	});
});

corviApp.controller('studyFinishedController', function($scope, $routeParams, $location, $log, Boxes) {
	
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
	}
	
	$scope.box = {};
	Boxes.get(id, function(data) {
		$scope.box = data;
	}, function(res) {
		$location.path("/");
	});
	
	$scope.getBack = function() {
		$location.path("/");
	};
	
});

corviApp.controller('studyQuestionController', function($scope, $routeParams, $log, $location, Questions) {
	
	var boxID = parseInt($routeParams.box, 10);
	if (isNaN(boxID)) {
		$log.error("Invalid ID!");
	}

	loadNewQuestion(boxID);	
	
	$scope.showSolution = function() {
		$scope.answered = true;	
	};
	
	$scope.giveCorrectAnswer = function() {
		Questions.giveCorrectAnswer($scope.question.ID, function() {
			loadNewQuestion(boxID);	
		}, function() {
			$location.path("/");
		});
	};
	
	$scope.giveWrongAnswer = function() {
		Questions.giveWrongAnswer($scope.question.ID, function() {
			loadNewQuestion(boxID);	
		}, function() {
			$location.path("/");
		});
	};
	
	function loadNewQuestion(boxID) {
		$scope.answered = false;
		$scope.question = {};
		
		Questions.getQuestionToLearn(boxID, function(data) {
			$scope.question = data;
		}, function() {
			$location.path("/study/"+boxID+"/finished");
		}, function(res) {
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
	};
});

corviApp.controller('mainNavigationController', function($scope, $route) {
	$scope.$route = $route;
	
	$scope.studyBadge = 3;
});