var corviApp = angular.module('corviApp', ['ngRoute']);

corviApp.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		templateUrl: 'sites/study_boxes.html',
		controller: 'studyBoxController',
		navActive: 'study'
	})
	.when('/study/:box', {
		templateUrl: 'sites/study_question.html',
		controller: 'mainController',
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
	});
});

corviApp.controller('studyBoxController', function($scope, $http, $log) {
	$scope.categories = {};
	$http.get("/api/categories").then(function(res) {
		if(res.status === 200) {
			$scope.categories = res.data;
			$log.debug(res.data);
		}else{
			$log.error(res);
		}
		
		Object.keys($scope.categories).forEach(function(key) {
			$http.get("/api/category/" + $scope.categories[key].ID + "/boxes").then(function(res) {
				if(res.status == 200) {
					$scope.categories[key].Boxes = res.data;
				}
				$log.debug($scope.categories);
			});
		});
		
	});
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
});