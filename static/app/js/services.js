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
    .factory('gadgets', ['$rootScope', '$http', function($rootScope, $http) {
        return {
            get: function(callback, errback) {
                $http.get('/gadgets').success(function (data, status, headers, config) {
                    callback(data);
                }).error(function(data, status, headers, config) {
                    errback();
                });
            }
        }
    }])
    .factory('methods', ['$rootScope', '$http', function($rootScope, $http) {
        return {
            save: function(name, method) {
                var url, httpMethod, data
                if (method.id != undefined && method.id > 0) {
                    url = '/gadgets/' + name + '/methods/' + method.id.toString();
                    httpMethod = 'PUT'
                } else {
                    url = '/gadgets/' + name + '/methods';
                    httpMethod = 'POST'
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
            get: function(name, callback) {
                var url = '/gadgets/' + name + '/methods';
                $http.get(url).success(function (data, status, headers, config) {
                    callback(data);
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
                    data: JSON.stringify({username:username, password: password}),
                    headers: {'Content-Type': 'application/json'}
                }).success(function (data, status, headers, config) {
                    return true;
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
                    
                });
            }
        }
    }]);




