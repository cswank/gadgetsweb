'use strict';

angular.module('myApp.services', [])
    .value('version', '0.1')
    .factory('sockets', ['$rootScope', '$location', function($rootScope, $location) {
        var ws;
        var outWs;
        var subscribeCallbacks = [];
        var canWrite;
        
        function getWebsockets(gadget) {
            ws = {};
            var prot = "wss";
            if ($location.protocol() == "http") {
                prot = "ws";
            }
            var url = prot + "://" + $location.host() + "/api/socket/in?host=" + gadget;
            ws.input = new WebSocket(url);
            var url = prot + "://" + $location.host() + "/api/socket/out?host=" + gadget;
            ws.output = new WebSocket(url);
            return ws;
        }
        return {
            connect: function(gadget, errorCallback) {
                if(ws != undefined) {
                    ws.input.close();
                    ws.input = null;
                    ws.output.close();
                    ws.output = null;
                }
                ws = getWebsockets(gadget);
                ws.input.onopen = function() {
                };
                ws.input.onerror = function() {
                };
                ws.output.onopen = function() {
                    canWrite = true;
                };
                ws.output.onerror = function(data) {
                    console.log(data);
                };
                ws.input.onmessage = function(message) {
                    message = JSON.parse(message.data);
                    var event = message[0];
                    if (event == 'ping') {
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
                if (canWrite) {
                    ws.output.send(message);
                }
            },
            subscribe: function(callback) {
                subscribeCallbacks.push(callback);
            },
            close: function() {
                if (ws != undefined) {
                    ws.input.close();
                    ws.output.close();
                }
            }
        }
    }])
    .factory('gadgets', ['$http', '$location', function($http, $location) {
        return {
            get: function(callback, error) {
                $http.get('/api/gadgets').success(function (data, status, headers, config) {
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
                var url = '/api/history/gadgets/' + name + '/devices';
                $http.get(url).success(function (data, status, headers, config) {
                    callback(data);
                }).error(function(data, status, headers, config) {
                    console.log(data);
                });
            }
        }
    }])
    .factory('methods', ['$rootScope', '$http', function($rootScope, $http) {
        return {
            save: function(name, method) {
                var url, httpMethod, data
                if (method.id != undefined && method.id > 0) {
                    url = '/api/gadgets/' + name + '/methods/' + method.id.toString();
                    httpMethod = 'PUT';
                } else {
                    url = '/api/gadgets/' + name + '/methods';
                    httpMethod = 'POST';
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
                var url = '/api/gadgets/' + name + '/methods';
                $http.get(url).success(function (data, status, headers, config) {
                    callback(data);
                }).error(function() {
                    
                });
            },
            delete: function(name, method, callback) {
                var url = '/api/gadgets/' + name + '/methods/' + method.id.toString();
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
                    url: '/api/login',
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
                    url: '/api/logout',
                    method: "POST",
                    headers: {'Content-Type': 'application/json'}
                }).success(function (data, status, headers, config) {
                    callback();
                }).error(function (data, status, headers, config) {
                    
                });
            }
        }
    }]);




