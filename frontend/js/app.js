(function() {

var corviApp = angular.module('corviApp', ['ngRoute', 'chart.js', 'angularMoment', 'corviServices', 'corviCategoryControllers', 'corviBoxControllers', 'corviQuestionControllers']);

corviApp.run(function($window, Notify, Categories, Boxes, Questions, Settings, Stats) {
	Notify.connect();
	
	Categories.Refresh();
	Boxes.Refresh();
	Questions.Refresh();
	Settings.Refresh();
	Stats.Refresh();
	
	$window.onbeforeunload = function() {
		Notify.Destroy();
	};
});

corviApp.config(function($routeProvider, ChartJsProvider) {
	
	ChartJsProvider.setOptions({
		responsive: true,
		maintainAspectRatio: true
	});
	
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
		controller: 'statsController',
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

corviApp.controller('studyFinishedController', function($scope, $routeParams, $location, $log, $document, Boxes) {
	
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}
	
	$scope.boxID = id;
	$scope.box = Boxes.BoxesByID[id];
	
	$scope.keyPressed = function(event) {
		if (event.keyCode == 40) { // Down Arrow
			$scope.getBack();
		}
	};
	
	$document.on('keyup', $scope.keyPressed);
	$scope.$on('$destroy', function () {
		$document.off('keyup', $scope.keyPressed);
	});
	
	$scope.getBack = function() {
		$location.path("/");
		$scope.$apply();
	};
	
});

corviApp.controller('studyQuestionController', function($scope, $routeParams, $log, $location, $document, Questions) {
	
	$scope.answered = false;
	$scope.question = {};
	
	var boxID = parseInt($routeParams.box, 10);
	if (isNaN(boxID)) {
		$log.error("Invalid ID!");
		return
	}

	loadNewQuestion(boxID);	
	
	$scope.keyPressed = function(event) {
		if (event.keyCode == 40 && !$scope.answered) { // Down Arrow
			$scope.showSolution();
		}
		if (event.keyCode == 37 && $scope.answered) { // Left Arrow
			$scope.giveCorrectAnswer();
		}
		if (event.keyCode == 39 && $scope.answered) { // Right Arrow
			$scope.giveWrongAnswer();
		}
	};
	
	$scope.showSolution = function() {
		$scope.answered = true;	
		$scope.$apply();
	};
	
	$document.on('keyup', $scope.keyPressed);
	$scope.$on('$destroy', function () {
		$document.off('keyup', $scope.keyPressed);
	});
	
	$scope.giveCorrectAnswer = function() {
		if (!$scope.answered) {
			return
		}
		Questions.giveCorrectAnswer($scope.question.ID, function() {
			loadNewQuestion(boxID);	
		}, function() {
			$location.path("/");
		});
	};
	
	$scope.giveWrongAnswer = function() {
		if (!$scope.answered) {
			return
		}
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
	$scope.aboutWindow = null;
	$scope.openAboutWindow = function() {
		if ($scope.aboutWindow != null) {
			return
		}
		$scope.aboutWindow = new BrowserWindow({ 'width': 320, 'height': 190, 'show': false, 'title': 'About Corvi', 'titleBarStyle': 'hidden-inset', 'resizable': false, 'minimizable': false, 'maximizable': false });
		$scope.aboutWindow.loadURL('http://localhost:8080/app/about.html');
		$scope.aboutWindow.on('closed', function() {
  			$scope.aboutWindow = null;
		});
		$scope.aboutWindow.show();
	};
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

corviApp.controller('statsController', function($scope, $log, Stats) {

	$scope.range = Stats.Range;
	$scope.stats = Stats.Stats;
	$scope.rangeStr = "today";
	$scope.weekdays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
	$scope.monthDays = [];
	for (var i = 1; i <= 31; i++) {
		$scope.monthDays.push(i);
	}
	$scope.months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
	$scope.series = ['Answers'];
	
	$scope.isActive = function(range) {
		return range == $scope.rangeStr;
	};
	
	$scope.setRange = function(range) {
		$scope.rangeStr = range;
		switch(range) {
			case "today":
				var from = new Date(); // Today
				from.setHours(0,0,0,0,0);
				
				var to = new Date(+new Date() + 86400000); // Tomorrow
				to.setHours(0,0,0,0,0);
				
				Stats.SetRange(from, to);
				break;
				
			case  "week":
				var today = new Date();
  				var day = today.getDay();
      			var diff = today.getDate() - day + (day == 0 ? -6 : 1);
				var from = new Date(today.setDate(diff));
				from.setHours(0,0,0,0,0);

				var to = new Date(+new Date() + 86400000);
				to.setHours(0,0,0,0,0);
				
				Stats.SetRange(from, to);
				break;
				
			case "month":
				var from = new Date();
				from.setDate(1);
				from.setHours(0,0,0,0,0);
			
				var to = new Date(+new Date() + 86400000); // Tomorrow
				to.setHours(0,0,0,0,0);
				
				Stats.SetRange(from, to);
				break;
				
			case "year":
				var from = new Date();
				from.setDate(1);
				from.setMonth(0);
				from.setHours(0,0,0,0,0);
			
				var to = new Date(+new Date() + 86400000); // Tomorrow
				to.setHours(0,0,0,0,0);
				
				Stats.SetRange(from, to);
				break;
				
			case "all":
				var from = new Date(0);
				from.setHours(0,0,0,0,0);
			
				var to = new Date(+new Date() + 86400000); // Tomorrow
				to.setHours(0,0,0,0,0);
				
				Stats.SetRange(from, to);
				break;
		}
	};
	
});

corviApp.filter('percentage', ['$filter', function($filter) {
	return function(input, decimals) {
		return $filter('number')(input * 100, decimals) + '%';
	};
}]);

})();