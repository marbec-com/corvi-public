(function() {

var corviApp = angular.module('corviCategoryControllers', ['ngRoute', 'corviServices']);

corviApp.controller('categoryEditController', function($scope, $routeParams, $log, $location, Categories) {
	var id = parseInt($routeParams.category, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
	}
	
	$scope.category = angular.copy(Categories.CategoriesByID[id]);
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.save = function() {
		Categories.Update(id, $scope.category, function(data) {
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

corviApp.controller('categoryAddController', function($scope, $log, $location, Categories) {
	$scope.category = {};
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.save = function() {
		Categories.Add($scope.category, function(data) {
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

})();