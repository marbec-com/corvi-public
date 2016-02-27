(function() {

var corviApp = angular.module('corviBoxControllers', ['ngRoute', 'corviServices']);

corviApp.controller('boxEditController', function($scope, $routeParams, $log, $location, Boxes, Categories) {
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}
	
	$scope.box = angular.copy(Boxes.BoxesByID[id]);
	$scope.categories = Categories.CategoriesByID;
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.save = function() {
		Boxes.Update(id, $scope.box, function(data) {
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

corviApp.controller('boxAddController', function($scope, $log, $location, Boxes, Categories) {
	$scope.box = {};
	$scope.categories = Categories.CategoriesByID;
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.save = function() {
		Boxes.Add($scope.box, function(data) {
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

})();