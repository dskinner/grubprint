angular.module("services", ["ngResource"]).
	factory("Food", function($resource) {
		var Food = $resource("/foodQuery?q=:q");
		return Food;
	}).
	factory("NutrientData", function($resource) {
		var NutrientData = $resource("/nutrientDataQuery?id=:id");
		return NutrientData;
	});
