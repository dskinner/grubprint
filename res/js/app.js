angular.module("food", ["ui.bootstrap", "ngRoute", "ngResource", "ngGrid", "services", "filters"]).
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
				templateUrl: "/html/food.html"
			}).
			otherwise({redirectTo: "/"});
	});
