angular.module("services", ["ngResource"]).
	factory("Food", function($resource) {
		var Food = $resource("/foods/:q");
		return Food;
	}).
	factory("NutrientData", function($resource) {
		var NutrientData = $resource("/nutrients/:id");
		return NutrientData;
	}).
	factory("Weight", function($resource) {
		var Weight = $resource("/weights/:id");
		return Weight;
	});
