angular.module("services", ["ngResource"]).
	factory("Food", function($resource) {
		var Food = $resource("/foodQuery?q=:projectId");
		return Project;
	});
