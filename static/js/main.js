var app = angular.module('steel', ['ngRoute', 'ui.codemirror']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.
        when('/classes/:class_id/assignments/:id?', {
            templateUrl: 'templates/assignments.html',
            controller: 'AssignmentsCtrl'
        }).
        when('/classes/:class_id/assignments/:id/:team', {
            templateUrl: 'templates/personal_assignment.html',
            controller: 'PersonalAssignmentCtrl'
        }).
        when('/class/add', { // move this to /classes/
            templateUrl: 'templates/add_class.html',
            controller: 'AddClassCtrl'
        }).
        when('/classes/:class_id/sandbox', {
            templateUrl: 'templates/sandbox.html',
            controller: 'SandboxCtrl'
        }).
        when('/classes/:class_id/admin', {
            templateUrl: 'templates/admin.html',
            controller: 'AdminCtrl'
        }).
        when('/auth/:method', {
            templateUrl: 'templates/auth.html',
            controller: 'AuthCtrl'
        }).
        when('/login', {
            redirectTo: '/auth/login'
        }).
        when('/register', {
            redirectTo: '/auth/register'
        }).
        when('/classes/:class_id/', {
            templateUrl: 'templates/feed.html',
            controller: 'FeedCtrl'
        }).
        when('/', {
            templateUrl: 'templates/hello.html'
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

app.factory('user', ['$http', '$location', '$rootScope', function($http, $location, $rootScope) {
    var u = {
        username: undefined,
        token: localStorage.token,
        admin: false,
        classes: [],
        classID: undefined,
        setToken: function(t) {
            console.log('changed token to:', t);
            u.token = t;
            localStorage.token = t;
            
            $rootScope.$broadcast('user.update');
        },
        auth: function(method, username, password, email)  {
            var url = '/users/' + method + '?username=' + username + '&password=' + password + '&email=' + email;
            console.log('url: ', url);

            $http.post(url).
                then(function(resp) {
                    // TODO error handling
                    if (resp.data.data.username !== undefined) {
                        console.log(resp.data);
                        $location.path('/auth/login');
                        return;
                    }

                    u.username = username;
                    u.setToken(resp.data.data);

                    $rootScope.$broadcast('user.update');

                    u.classes();

                    $location.path('/');
                });
        },
        loggedIn: function() {
            $http.get('/users?access_token=' + u.token).
                then(function(resp) {
                    if (resp.data.status.error) {
                        $location.path('/auth/login');
                        return;
                    }

                    console.log(resp.data);

                    u.username  = resp.data.data.username;

                    u.classes(); 

                    $rootScope.$broadcast('user.update');
                });
        },
        classes: function() {
            $http.get('/classes?access_token=' + u.token)
                .then(function(resp) {
                    if (resp.data.status.error) {
                        console.log("COULD NOT GET USERS");
                        return;
                    }

                    u.classes = resp.data.data;
                    $rootScope.$broadcast('user.update');
                });
        },
        setClass: function(i) {
            u.classID = i;
            u.currentStudent();
            $rootScope.$broadcast('user.update');
        },
        currentStudent: function() {
            $http.get('/classes/' + u.classID + '/students?access_token=' + u.token).
                then(function(resp) {
                    // TODO error handling
                    u.student = resp.data.data;
                    $rootScope.$broadcast('user.update');
                });
        }
    };

    u.loggedIn();

    return u;
}]);

app.controller('AddClassCtrl', ['$scope', '$http', '$location', '$routeParams', 'user', function($scope, $http, $location, $routeParams, user) {
    user.setClass($routeParams.class_id);
    $scope.go = function() {
        var name = $scope.name;
        var description = $scope.description;

        if (!name || !description) {
            return;
        }

        $http.post('/classes?access_token=' + user.token + '&name=' + name + '&description=' + description)
            .then(function(resp) {
                if (resp.data.status.error) {
                    console.log("COULD NOT CREATE NEW CLASSES");
                    console.log(resp);

                    return;
                }

                $location.path('/classes/' + resp.data.data.id);
            });
    }

}]);

app.controller('AdminCtrl', ['$scope', '$http', '$location', '$routeParams', 'user', function($scope, $http, $location, $routeParams, user) {
    user.setClass($routeParams.class_id);
    if (!user.admin) {
        $location.path('/');
        return;
    }
}]);

app.controller('WholeCtrl', ['$rootScope', 'user', '$location', function($rootScope, user, $location) {
    $rootScope.loggedIn = false;
    $rootScope.user = undefined;
    $rootScope.inClass = false;
    $rootScope.student = undefined;

    $rootScope.$on('user.update', function(evt) {
        $rootScope.user = user;
        $rootScope.loggedIn = true;
        $rootScope.inClass = user.classID !== undefined; // true;
        $rootScope.student = user.student;
        console.log("stud");
        console.log($rootScope.student);
    });
}]);

app.controller('AuthCtrl', ['$scope', '$routeParams', '$http', '$location', 'user', function($scope, $routeParams, $http, $location, user) {
    if ($routeParams.method != 'login' && $routeParams.method != 'register') {
        $location.path('/');
        return;
    }

    $scope.current = $routeParams.method;
    $scope.other = $scope.current == 'login' ? 'register' : 'login';

    $scope.go = function() {
        console.log($scope.username, $scope.password, $scope.email);
        user.auth($scope.current, $scope.username, $scope.password, $scope.email);
    };
}]);

app.controller('SandboxCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);
    $scope.editorOpts = {
        lineWrapping : true,
        lineNumbers: true,
        theme: 'elegant',
        mode: 'javascript',
        value: "// do some magic?"
    };
}]);

app.controller('PersonalAssignmentCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);
}]);

app.controller('AssignmentsCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);
    var id = parseInt($routeParams.id) || 1;

    $http.get('/classes/' + user.classID + '/assignments/due?access_token=' + user.token)
        .then(function(resp) {
            $scope.tasks = resp.data.data;
            for (var i = 0; i < $scope.tasks.length; i++) {
                var task = $scope.tasks[i]
                if (task.id == id) {
                    $scope.selected = task;
                    break;
                }
            }
        });

    $scope.tasks = [];
    $scope.selected = {};
}]);

app.controller('FeedCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);
    $http.get('/classes/' + user.classID + '/assignments/due?access_token=' + user.token)
        .then(function(resp) {
            $scope.dueTasks = resp.data.data;
            console.log(resp);
        });

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
