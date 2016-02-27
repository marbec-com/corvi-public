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
	
	$scope.questions = Questions.QuestionsByBoxID[id] || [];
});

})();