var app = angular.module('steel', ['ngRoute', 'ui.codemirror', 'hc.marked']);

app.config(['markedProvider', function(markedProvider) {
    markedProvider.setOptions({gfm: true});
}]);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.
        when('/classes/:class_id/workshops/create/:workshop_id?', {
            templateUrl: 'templates/create_workshop.html',
            controller: 'CreateWorkshopCtrl'
        }).
        when('/classes/:class_id/workshops/:id?', {
            templateUrl: 'templates/workshops.html',
            controller: 'WorkshopsCtrl'
        }).
        when('/classes/:class_id/workshops/:workshop_id/pages', {
            templateUrl: 'templates/view_workshop.html',
            controller: 'ViewWorkshopCtrl'
        }).
        when('/classes/:class_id/workshops/:id/:team', {
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
            templateUrl: 'templates/hello.html',
            controller: 'HelloCtrl'
        }).
        otherwise({
            redirectTo: '/'
        });
}]);

app.filter('titlecase', function() {
    return function(input) {
        var words = input.split(' ');
        for (i in words) {
            console.log(words[i]);
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
        broadcast: function() {
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
            for (var n = 0; n < u.classes.length; n++) {
                if (u.classes[n].id == i) {
                    u.class = u.classes[n];
                    break;
                }
            }

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

app.controller('ViewWorkshopCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);

    $scope.pages = [];
    $scope.selected = {};

    $http.get('/classes/'+ user.classID + '/workshops/' + $routeParams.workshop_id + '/pages?access_token=' + user.token).
        then(function(resp) {
            console.log(resp);
            $scope.pages = resp.data.data;
            $scope.selected = $scope.pages[0];
        });

    $scope.view = function(page) {
        $scope.selected = page;
    };
}]);

app.controller('CreateWorkshopCtrl', ['$scope', '$http', '$location', '$routeParams', 'user', function($scope, $http, $location, $routeParams, user) {
    user.setClass($routeParams.class_id);

    $scope.editorOpts = {
        lineNumbers: true,
        theme: 'ambiance',
        mode: 'markdown'
    }

    $scope.workshop = {};
    $scope.canCreatePages = false;
    $scope.edit = false;

    $scope.pages = [];

    if ($routeParams.workshop_id !== undefined) {
        $scope.edit = true;

        $http.get('/classes/' + $routeParams.class_id + '/workshops/' + $routeParams.workshop_id + '?access_token=' + user.token).
            then(function(resp) {
                console.log(resp.data);
                $scope.workshop = resp.data.data;
            });

        $http.get('/classes/' + $routeParams.class_id + '/workshops/' + $routeParams.workshop_id + '/pages?access_token=' + user.token).
            then(function(resp) {
                console.log(resp.data);
                $scope.pages = resp.data.data;
            });
    }

    $scope.create = function() {
        $http.post('/classes/' + $routeParams.class_id + '/workshops?access_token=' + user.token + '&name=' + $scope.workshop.name + '&description=' + $scope.workshop.description + '&explanation=' + $scope.workshop.explanation).
            then(function(resp) {
                console.log(resp.data);
                $scope.edit = true;
                $scope.canCreatePages = true;
                $scope.workshop = resp.data.data;
            });
    }

    $scope.tempWorkshop = function() {
        $scope.pages.push({created: false});
    }

    $scope.delete = function() {
        $http.delete('/classes/' + $routeParams.class_id + '/workshops/'+ $scope.workshop.id +'?access_token=' + user.token).
            then(function(resp) {
                $location.path('/classes/' + $routeParams.class_id + '/workshops');
            });
    };

    $scope.addPage = function() {
        console.log($scope.pages);
        var page = $scope.pages[$scope.pages.length - 1];
        $http.post('/classes/' + user.classID + '/workshops/' + $scope.workshop.id + '/pages?access_token=' + user.token + '&title=' + page.title +  '&content=' + page.contents).
            then(function(resp) {
                console.log(resp.data);
                $scope.pages[$scope.pages.length - 1].created = true;
            });
    }

    $scope.updatePage = function(page) {
        console.log("updating", page);
        console.log(page.contents);
        $http({
            method: 'POST',
            url: '/classes/' + user.classID + '/workshops/' + $scope.workshop.id + '/pages/' + page.id + '/edit',
            params: {access_token: user.token, contents: page.contents, title: page.title}
        }).
            then(function(resp) {
                console.log(resp.data);
            });
    }

    $scope.updateWorkshop = function() {
            $http({
                method: 'PUT',
                url:'/classes/' + user.classID + '/workshops/' + $scope.workshop.id,
                params: {access_token: user.token, name: $scope.workshop.name, explanation: $scope.workshop.explanation, description: $scope.workshop.description}
            }).
            then(function(resp) {
                console.log(resp);
            });
            /*$http.put('/classes/' + user.classID + '/workshops/' + $scope.workshop.id, {params: {access_token: user.token, name: $scope.workshop.name, explanation: $scope.workshop.explanation, description: $scope.workshop.description}}).
                then(function(resp) {
                    console.log(resp);
                });*/
    }

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
    if ($scope.student == undefined) {
        $location.path('/');
    }
    
    if ($scope.student.permission != 1) {
        $location.path('/');
        return;
    }

    $scope.invite = function() {
        console.log("inviting");
        $http.post('/classes/' + user.classID + '/admin/students?access_token=' + user.token + '&user=' + $scope.invitee).
            then(function(resp) {
                console.log(resp.data.message);
            });
    }

    $scope.image = function() {
        console.log("Adding image");

        $http.post('/classes/' + user.classID + '/image?access_token=' + user.token + '&image_url=' + $scope.image_url).
            then(function(resp) {
                console.log(resp)
            });

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

app.controller('WorkshopsCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);

    $scope.tasks = [];
    $scope.selected = {};

    var id = $routeParams.id
    if (id !== undefined) {
        var id = parseInt($routeParams.id);
    }

    $http.get('/classes/' + user.classID + '/workshops?access_token=' + user.token)
        .then(function(resp) {
            if (resp.data.data == null) {
                $scope.selected = {
                    name: "Lorem Ipsum",
                    description: "How do I like my pie? Shaken not stirred.",
                    explanation: "a^2 = b^2 + c^2 - 2*b*c*cosA"
                }

                return;
            }

            $scope.tasks = resp.data.data;
            for (var i = 0; i < $scope.tasks.length; i++) {
                var task = $scope.tasks[i]
                if (task.id == id || id === undefined) {
                    $scope.selected = task;
                    break;
                }
            }
        });

    $scope.view = function(task) {
        $scope.selected = task;
    }
}]);

app.controller('FeedCtrl', ['$scope', '$http', '$routeParams', 'user', function($scope, $http, $routeParams, user) {
    user.setClass($routeParams.class_id);
    $http.get('/classes/' + user.classID + '/workshops?access_token=' + user.token)
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

app.controller('HelloCtrl', ['$scope', '$http', 'user', function($scope, $http, user) {
    user.student = undefined;
    user.classID = undefined;
    user.broadcast();

    $scope.classes = [];
    $http.get('/classes?access_token=' + user.token)
        .then(function(res) {
            $scope.chunkedClasses = chunk(res.data.data, 3);
            console.log($scope.chunkedClasses);
        });
}]);

chunk = function(array, chunkSize) {
    return [].concat.apply([],
            array.map(function(elem,i) {
                return i%chunkSize ? [] : [array.slice(i,i+chunkSize)];
            })
            );
}
