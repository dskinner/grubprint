function FoodCtrl($scope, Food) {
	$scope.foods = Food.get();
}
