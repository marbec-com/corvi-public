(function() {

var corviApp = angular.module('corviApp', ['ngRoute', 'angularMoment', 'corviServices', 'corviCategoryControllers', 'corviBoxControllers', 'corviQuestionControllers']);

corviApp.run(function($window, Notify, Categories, Boxes, Questions, Settings) {
	Notify.connect();
	
	Categories.Refresh();
	Boxes.Refresh();
	Questions.Refresh();
	Settings.Refresh();
	
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
	.when('/manage/category/add', {
		templateUrl: 'sites/category_add.html',
		controller: 'categoryAddController',
		navActive: 'manage'
	})
	.when('/manage/category/:category/edit', {
		templateUrl: 'sites/category_edit.html',
		controller: 'categoryEditController',
		navActive: 'manage'
	})
	.when('/manage/category/:category/delete', {
		templateUrl: 'sites/category_delete.html',
		controller: 'categoryDeleteController',
		navActive: 'manage'
	})
	.when('/manage/box/add', {
		templateUrl: 'sites/box_add.html',
		controller: 'boxAddController',
		navActive: 'manage'
	})
	.when('/manage/box/:box/edit', {
		templateUrl: 'sites/box_edit.html',
		controller: 'boxEditController',
		navActive: 'manage'
	})
	.when('/manage/box/:box/delete', {
		templateUrl: 'sites/box_delete.html',
		controller: 'boxDeleteController',
		navActive: 'manage'
	})
	.when('/manage/box/:box/questions', {
		templateUrl: 'sites/questions.html',
		controller: 'questionsController',
		navActive: 'manage'
	})
	.when('/manage/box/:box/questions/add', {
		templateUrl: 'sites/question_add.html',
		controller: 'questionAddController',
		navActive: 'manage'
	})
	.when('/manage/question/:question/edit', {
		templateUrl: 'sites/question_edit.html',
		controller: 'questionEditController',
		navActive: 'manage'
	})
	.when('/manage/question/:question/delete', {
		templateUrl: 'sites/question_delete.html',
		controller: 'questionDeleteController',
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
		controller: 'settingsEditController',
		navActive: 'settings'
	}).otherwise({
		redirectTo: '/'
	});
});

corviApp.controller('studyBoxController', function($scope, $log, Categories, Boxes, Settings) {
	$scope.categories = Categories.CategoriesAll;
	$scope.boxesByCatID = Boxes.BoxesByCatID;
});

corviApp.controller('studyFinishedController', function($scope, $routeParams, $location, $log, Boxes) {
	
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
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
		return
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
});

corviApp.directive('mainNavigation', function() {
	return {
		retrict: 'E',
		templateUrl: 'partials/mainNavigation.html',
		controller: 'mainNavigationController'
	};
});

corviApp.controller('mainNavigationController', function($scope, $route, Boxes) {
	$scope.$route = $route;
	$scope.boxMetadata = Boxes.Metadata;
});

corviApp.controller('settingsEditController', function($scope, $log, $location, Settings) {
	
	$scope.settings = angular.copy(Settings.Settings);
	$scope.error = "";
	
	$scope.save = function() {
		Settings.Update($scope.settings, function(data) {
			$location.path("/settings");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

})();