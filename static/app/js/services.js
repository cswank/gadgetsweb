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
                    errorCallback();
                };
                ws.onmessage = function(message) {
                    message = JSON.parse(message.data);
                    var event = message[0];
                    var payload = JSON.parse(message[1]);
                    for (i in subscribeCallbacks) {
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
                ws.close();
            }
        }
    }])
    .factory('gadgets', ['$rootScope', '$http', function($rootScope, $http) {
        return {
            get: function(callback) {
                $http.get('/gadgets').success(function (data, status, headers, config) {
                    
                    callback(data);
                }).error(function(data, status, headers, config) {
                    
                    $rootScope.emit("login");
                });
            }
        }
    }])
    .factory('auth', ['$http', function($http) {
        return {
            login: function(username, password) {
                $http({
                    url: '/login',
                    method: "POST",
                    data: JSON.stringify({username:$scope.username, password: $scope.password}),
                    headers: {'Content-Type': 'application/json'}
                }).success(function (data, status, headers, config) {
                    return true;
                }).error(function (data, status, headers, config) {
                    return false;
                });
            }
        }
    }]);


