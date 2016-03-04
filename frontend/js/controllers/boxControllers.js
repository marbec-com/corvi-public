(function() {

var corviApp = angular.module('corviBoxControllers', ['ngRoute', 'corviServices']);

corviApp.controller('boxEditController', function($scope, $routeParams, $log, $location, Boxes, Categories) {
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}
	
	$scope.box = angular.copy(Boxes.BoxesByID[id]);
	$scope.categories = Categories.CategoriesAll;
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.save = function() {
		Boxes.Update(id, $scope.box, function() {
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

corviApp.controller('boxAddController', function($scope, $log, $location, Boxes, Categories, Questions) {
	$scope.box = {};
	$scope.categories = Categories.CategoriesAll;
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.save = function() {
		Boxes.Add($scope.box, function(box) {
			Questions.RefreshBox(box.ID);
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

corviApp.controller('boxDeleteController', function($scope, $routeParams, $log, $location, Boxes) {
	var id = parseInt($routeParams.box, 10);
	if (isNaN(id)) {
		$log.error("Invalid ID!");
		return
	}
	
	$scope.box = angular.copy(Boxes.BoxesByID[id]);
	$scope.error = "";
	
	$scope.back = function() {
		$location.path("/manage");
	};
	
	$scope.submit = function() {
		Boxes.Delete(id, function() {
			$location.path("/manage");
		}, function(err) {
			$scope.error = err;
		});
	};
	
});

})();