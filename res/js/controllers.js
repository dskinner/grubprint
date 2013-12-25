function FoodCtrl($scope, Food, NutrientData) {
	$scope.update = function() {
		$scope.foods = Food.query({q: $scope.search});
		$scope.totalServerItems = $scope.foods.length;
		$scope.updateNutrients();
	};

	$scope.updateNutrients = function() {
		if ($scope.selections.length === 0) return;
		$scope.nutrients = NutrientData.query({id: $scope.selections[0].Id});
	};

	$scope.foods = [];
	$scope.selections = [];
	$scope.nutrients = [];

	$scope.$watch("selections", $scope.updateNutrients);

	$scope.filterOptions = {
		filterText: "Beans",
		useExternalFilter: false
	};

	$scope.gridOptions = {
		data: "foods",
		showGroupPanel: true,
		columnDefs: [
			{field: "LongDesc", displayName: "Name"},
			{field: "NitrogenFactor"},
			{field: "ProteinFactor"},
			{field: "FatFactor"},
			{field: "CarbohydrateFactor"}
		],
		showFooter: true,
		filterOptions: $scope.filterOptions,
		selectedItems: $scope.selections,
		multiSelect: false,
		afterSelectionChange: function (rowItem, event) {
			$scope.updateNutrients();
		},
		enableColumnResize: true,
	};

	$scope.update();
}
