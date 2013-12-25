angular.module("food", ["services", "filters", "ui", "ui.bootstrap"]).
	config(function($httpProvider) {
		$httpProvider.defaults.transformRequest = function(data) {
			if (data !== undefined) {
				return $.param(data);
			}
		};
		$httpProvider.defaults.headers.post["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8";
	}).
	config(function($routeProvider) {
		$routeProvider.
			when("/", {
				controller: FoodCtrl,
				templateUrl: "/partial/food.html"
			}).
			otherwise({redirectTo: "/"});
	});
