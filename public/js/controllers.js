function FoodCtrl($scope, Food, NutrientData) {

	$scope.update = function() {
		$scope.foods = Food.query({q: $scope.search});
		$scope.totalServerItems = $scope.foods.length;
		$scope.updateNutrients();
	};

	$scope.updateNutrients = function() {
		if ($scope.selections.length === 0) return;
		$scope.nutrients = NutrientData.get({id: $scope.selections[0].Id});
	};

	$scope.foods = [];
	$scope.selections = [];
	$scope.nutrients = [];

	$scope.gridOptions = {
		data: "foods",
		showGroupPanel: false,
		columnDefs: [
			{field: "LongDesc", displayName: "Name"}
		],
		selectedItems: $scope.selections,
		multiSelect: false,
		enableColumnResize: true,
		showFooter: true,
		afterSelectionChange: function (rowItem, event) {
			$scope.updateNutrients();
		},
	};

	$scope.update();
}
