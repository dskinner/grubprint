function FoodCtrl($scope, Food, Weight, NutrientData) {

	$scope.update = function() {
		$scope.foods = Food.query({q: $scope.search});
		$scope.totalServerItems = $scope.foods.length;
		$scope.updateSelection();
	};

	$scope.updateSelection = function() {
		if ($scope.selections.length === 0) return;
		$scope.weights = Weight.query({id: $scope.selections[0].Id}, function() {
			$scope.nutrients = NutrientData.get({id: $scope.selections[0].Id});
		});
	};

	$scope.foods = [];
	$scope.selections = [];
	$scope.weights = [];
	$scope.nutrients = [];

	$scope.weight = 100;

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
			$scope.updateSelection();
		},
	};

	$scope.update();
}
