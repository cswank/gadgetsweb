'use strict';

angular.module('myApp.services', [])
    .value('version', '0.1')
    .factory('sockets', ['$rootScope', function($rootScope) {
        var ws;
        var subscribeCallbacks = [];
        return {
            connect: function(gadget, errorCallback) {
                if(ws) {
                    ws.close();
                    ws = null;
                }
                ws = new WebSocket("wss://gadgets.dyndns-ip.com/socket?host=" + gadget);
                ws.onopen = function() {
                };
                ws.onerror = function() {
                };
                ws.onmessage = function(message) {
                    message = JSON.parse(message.data);
                    var event = message[0];
                    if (event == 'ping') {
                        console.log("ping")
                        return;
                    }
                    var payload = JSON.parse(message[1]);
                    for (var i in subscribeCallbacks) {
                        var cb = subscribeCallbacks[i];
                        cb(event, payload);
                    }
                };
            },
            send: function(message) {
                ws.send(message);
            },
            subscribe: function(callback) {
                subscribeCallbacks.push(callback);
            },
            close: function() {
                if (ws != undefined) {
                    ws.close();
                }
            }
        }
    }])
    .factory('gadgets', ['$http', '$location', function($http, $location) {
        return {
            get: function(callback, error) {
                $http.get('/gadgets').success(function (data, status, headers, config) {
                    callback(data);
                }).error(function(data, status, headers, config) {
                    error();
                });
            }
        }
    }])
    .factory('history', ['$http', function($http) {
        return {
            getDevices: function(name, callback) {
                var url = '/history/gadgets/' + name + '/devices';
                console.log(url);
                $http.get(url).success(function (data, status, headers, config) {
                    callback(data);
                }).error(function(data, status, headers, config) {
                    console.log(data);
                });
            }
        }
<<<<<<< HEAD
    })
    .factory('methods', function($http, $rootScope, $modal) {
        function getMethods() {
            var url = '/gadgets/' + $scope.gadget.name + '/methods';
            $http.get(url).success(function (data, status, headers, config) {
                $scope.showMethods = true;
                $scope.methods = [$scope.method];
                for (var i in data.methods) {
                    var rawMethod = data.methods[i];
                    $scope.methods.push(rawMethod);
                }
            });
        }
        return {
            runMethod: function(method) {
                var msg = {event: 'method', message: {type: 'method', method:method}};
                socket.send(JSON.stringify(msg));
            },
            clearMethod: function() {
                var msg = {event: 'command', message: {type: 'command', body:'clear method'}};
                socket.send(JSON.stringify(msg));
            },
            saveMethod: function(method, gadget) {
                var url, httpMethod, data
                if (method.id != undefined && method.id > 0) {
                    url = '/gadgets/' + gadget.name + '/methods/' + method.id.toString();
                    httpMethod = 'PUT'
                } else {
                    url = '/gadgets/' + gadget.name + '/methods';
                    httpMethod = 'POST'
=======
    }])
    .factory('methods', ['$rootScope', '$http', function($rootScope, $http) {
        return {
            save: function(name, method) {
                var url, httpMethod, data
                if (method.id != undefined && method.id > 0) {
                    url = '/gadgets/' + name + '/methods/' + method.id.toString();
                    httpMethod = 'PUT';
                } else {
                    url = '/gadgets/' + name + '/methods';
                    httpMethod = 'POST';
>>>>>>> dc3a151856979f508ef20a8c95b18e603d0412f8
                }
                $http({
                    url: url,
                    method: httpMethod,
                    data: JSON.stringify(method),
                    headers: {'Content-Type': 'application/json'}
                }).success(function (data, status, headers, config) {
                    
                }).error(function (data, status, headers, config) {
                    console.log("error saving method");
                });
            },
<<<<<<< HEAD
            addMethod: function(method) {
                var dlg = $modal.open({
                    templateUrl: '/dialogs/method.html',
                    controller: MethodCtrl,
                    resolve: {
                        method: function () {
                            return method;
                        }
                    }
                });
                dlg.result.then(function(method) {
                    saveMethod();
                    return method;
                } ,function(){
                    
                });
            },
            getRecipe: function() {
                var dlg = $modal.open({
                    templateUrl: '/dialogs/recipe.html?c=' + new Date().getTime(),
                    controller: RecipeCtrl,
                });
                dlg.result.then(function(recipe) {
                    var url = '/recipes/' + recipe.name + '?grainTemperature=' + recipe.grainTemperature;
                    $http.get(url).success(function (data, status, headers, config) {
                        return data;
                        $scope.method = data;
                    });
                } ,function() {
=======
            get: function(name, callback) {
                var url = '/gadgets/' + name + '/methods';
                $http.get(url).success(function (data, status, headers, config) {
                    callback(data);
                }).error(function() {
                    
                });
            },
            delete: function(name, method, callback) {
                var url = '/gadgets/' + name + '/methods/' + method.id.toString();
                $http.delete(url).success(function (data, status, headers, config) {
                    callback(data);
                }).error(function() {
                    
                });
            }
        }
    }])
    .factory('auth', ['$http', function($http) {
        return {
            login: function(username, password, callback) {
                $http({
                    url: '/login',
                    method: "POST",
                    data: JSON.stringify({username:username, password: password}),
                    headers: {'Content-Type': 'application/json'}
                }).success(function (data, status, headers, config) {
                    callback();
                }).error(function (data, status, headers, config) {
                    return false;
                });
            },
            logout: function(callback) {
                $http({
                    url: '/logout',
                    method: "POST",
                    headers: {'Content-Type': 'application/json'}
                }).success(function (data, status, headers, config) {
                    callback();
                }).error(function (data, status, headers, config) {
>>>>>>> dc3a151856979f508ef20a8c95b18e603d0412f8
                    
                });
            }
        }
<<<<<<< HEAD
    });
=======
    }]);




>>>>>>> dc3a151856979f508ef20a8c95b18e603d0412f8
