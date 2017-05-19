var myApp = angular.module('myApp1', []);
myApp.controller('OptionsController', ['$scope', function($scope) {
	$scope.idd = 333;
	$scope.idStart = function(){
		$scope.idd = 0;
	};
	$scope.idTmp = function(){
		$scope.idd = $scope.idd+1;
		return $scope.idd;
	};
    $scope.option = [
		{name:"Temperatura", items:""},
		{name:"Riego automático", items:""},
		{name:"Ritmo cardiaco", items:""}
	];
}]);