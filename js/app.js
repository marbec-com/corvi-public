var corviApp = angular.module('corviApp', ['ngRoute', 'corviServices']);

corviApp.run(function($window, Notify, Categories, Boxes, Questions) {
	Notify.connect();
	
	Categories.Refresh();
	Boxes.Refresh();
	Questions.Refresh();
	
	$window.onbeforeunload = function() {
		Notify.Destroy();
	};
});

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
		controller: 'manageBoxesController',
		navActive: 'manage'
	})
	.when('/manage/addBox', {
		templateUrl: 'sites/manage.html',
		controller: 'mainController',
		navActive: 'manage'
	})
	.when('/manage/addCategory', {
		templateUrl: 'sites/manage.html',
		controller: 'mainController',
		navActive: 'manage'
	})
	.when('/manage/:box/edit', {
		templateUrl: 'sites/manage.html',
		controller: 'mainController',
		navActive: 'manage'
	})
	.when('/manage/:box/questions', {
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

corviApp.controller('studyBoxController', function($scope, $log, Categories, Boxes) {
	$scope.categories = Categories.CategoriesAll;
	$scope.boxesByCatID = Boxes.BoxesByCatID;
	
	//$log.debug($scope.categories);
});

corviApp.controller('studyFinishedController', function($scope, $routeParams, $location, $log, Boxes) {
	
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
	}
	
	$scope.boxID = id;
	$scope.box = Boxes.BoxesByID[id];
	
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
		
		Questions.GetQuestionToLearn(boxID, function(data) {
			$scope.question = data;
		}, function() {
			$location.path("/study/"+boxID+"/finished");
		}, function(res) {
			$location.path("/");
		});
	};
	
});

corviApp.controller('mainController', function($scope) {
});

corviApp.controller('manageBoxesController', function($scope, $log, Categories, Boxes) {
	$scope.categories = Categories.CategoriesAll;
	$scope.boxesByCatID = Boxes.BoxesByCatID;
	
	/* $scope.renameCategoryModal = function(catKey) {
		$(document).foundation();
		$('#categoryRenameModal').foundation('open');
		$scope.form.catKey = catKey;
		$scope.form.category = angular.copy($scope.categories[catKey]);
	}; 
	
	$scope.renameCategoryReset = function(form) {
		// Reset form elements
		if (form) {
			form.$setPristine();
			form.$setUntouched();
		}
		// Close and reset data
		$('#categoryRenameModal').foundation('close');
		$scope.form.catKey = -1;
		$scope.form.category = {};
	};
	
	$scope.renameCategorySave = function(form) {
		$log.debug(form);
		
		// Save
		$scope.categories[form.catKey] = angular.copy($scope.form.category);
		var copy = angular.copy($scope.form.category);
		Categories.update(copy, function() {});
		$scope.renameCategoryReset(form);
	};
	
	*/
		
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