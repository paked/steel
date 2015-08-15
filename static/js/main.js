var app = angular.module('steel', ['ngRoute', 'ui.codemirror']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.
        when('/assignments/:id?', {
            templateUrl: 'templates/assignments.html',
            controller: 'AssignmentsCtrl'
        }).
        when('/assignments/:id/:team', {
            templateUrl: 'templates/personal_assignment.html',
            controller: 'PersonalAssignmentCtrl'
        }).
        when('/sandbox', {
            templateUrl: 'templates/sandbox.html',
            controller: 'SandboxCtrl'
        }).
        when('/auth/:method', {
            templateUrl: 'templates/auth.html',
            controller: 'AuthCtrl'
        }).
        when('/', {
            templateUrl: 'templates/feed.html',
            controller: 'FeedCtrl'
        }).
        otherwise({
            redirectTo: '/'
        });
}]);

app.filter('titlecase', function() {
    return function(input) {
        var words = input.split(' ');
        for (i in words) {
            words[i] = words[i].toLowerCase();
            words[i] = words[i].charAt(0).toUpperCase() + words[i].slice(1)
        }

        return words.join(' ')
    };
});

app.controller('AuthCtrl', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location) {
    if ($routeParams.method != 'login' && $routeParams.method != 'register') {
        $location.path('/');
        return;
    }

    $scope.current = $routeParams.method;
    $scope.other = $scope.current == 'login' ? 'register' : 'login';

    $scope.go = function() {
        console.log($scope.username, $scope.password, $scope.email);
    };
}]);

app.controller('SandboxCtrl', ['$scope', '$http', function($scope, $http) {
    $scope.editorOpts = {
        lineWrapping : true,
        lineNumbers: true,
        theme: 'elegant',
        mode: 'javascript',
        value: "// do some magic?"
    };
}]);

app.controller('PersonalAssignmentCtrl', ['$scope', '$http', function($scope, $http) {

}]);

app.controller('AssignmentsCtrl', ['$scope', '$http', '$routeParams', function($scope, $http, $routeParams) {
    var index = parseInt($routeParams.id) || 1;

    $scope.tasks = [
        {
            "name": "Funny Strings",
            "description": "If the reverse of a character (a = z, b = y, c = x, etc.) is opposite iteself lorem ipsum xyz",
            "id": 1,
            "done": false,
            "time": "Due in 3 days"
        },
        {
            "name": "Reverse String",
            "description": "Given a string, how would you reverse it... quickly and easily",
            "id": 2,
            "done": false,
            "time": "Due in 3 days"
        },
        {
            "name": "Utopian Tree",
            "description": "Given a spec for a program, implement it using your knowledge of control structures",
            "id": 3,
            "done": true,
            "time": "Due in 4 days"
        },
        {
            "name": "Print",
            "description": "Make a few words appear in your terminal",
            "id": 4,
            "done": true,
            "time": "Due in 10 days"
        }
    ];

    $scope.selected = $scope.tasks[index - 1];
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
