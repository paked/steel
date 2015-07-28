var app = angular.module('steel', ['ngRoute']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.
        when('/assignments', {
            templateUrl: 'templates/assignments.html',
            controller: 'AssignmentsCtrl'
        }).
        when('/', {
            templateUrl: 'templates/feed.html',
            controller: 'FeedCtrl'
        }).
        otherwise({
            redirectTo: '/'
        });
}]);

app.controller('AssignmentsCtrl', ['$scope', '$http', function($scope, $http) {

}]);

app.controller('FeedCtrl', ['$scope', '$http', function($scope, $http) {
    $scope.dueTasks = [
        {
            "name": "Funny Strings",
            "id": 1,
            "done": false
        },
        {
            "name": "Reverse String",
            "id": 2,
            "done": false
        },
        {
            "name": "Utopian Tree",
            "id": 3,
            "done": true
        },
        {
            "name": "Print",
            "id": 4,
            "done": true
        }
    ];

    $scope.feedUpdates = [
        {
            "type": "assignment",
            "from": "Greg",
            "time": "20 minutes ago",
            "message": "Greg posted a new assignment"
        },
        {
            "type": "feedback",
            "from": "Greg",
            "time": "5 hours ago",
            "message": "Greg gave feedback on your \"Funny Strings\" task"
        },
        {
            "type": "like",
            "from": "Jimmy",
            "time": "10 hours ago",
            "message": "Jimmy liked your post"
        },
        {
            "type": "request",
            "from": "James",
            "time": "1 day ago",
            "message": "James requested to work with you on the Doomsday assignment"
        }
    ];
    
    $scope.allTasks = [
        
    ];
}]);
