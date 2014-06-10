'use strict';

angular.module('myApp.services', [])
    .value('version', '0.1')
    .factory('sockets', ['$rootScope', '$location', function($rootScope, $location) {
        var ws;
        var outWs;
        var subscribeCallbacks = [];
        
        function getWebsockets(gadget) {
            var prot = "wss";
            if ($location.protocol() == "http") {
                prot = "ws";
            }
            var url = prot + "://" + $location.host() + "/api/socket?host=" + gadget;
            ws = new WebSocket(url);
            return ws;
        }
        return {
            connect: function(gadget, errorCallback) {
                if(ws != undefined) {
                    ws.close();
                    ws = null;
                }
                ws = getWebsockets(gadget);
                ws.onopen = function() {
                };
                ws.onerror = function() {
                };
                ws.onmessage = function(message) {
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
    .factory('notes', ['$http', function($http) {
        return {
            save: function(name, note, callback) {
                var url = '/api/gadgets/' + name + '/notes';
                $http.post(url, note).success(function (data) {
                    callback();
                });
            },
            get: function(name, callback) {
                var url = '/api/gadgets/' + name + '/notes';
                $http.get(url).success(function (data) {
                    callback(data);
                });
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
                    console.log("error", data);
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




