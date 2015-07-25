var app = angular.module('steel', []);

app.controller('homeCtrl', ['$scope', function($scope){
    $scope.message = 'Hello!';
}]);
