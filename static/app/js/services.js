'use strict';

angular.module('myApp.services', [])
    .factory('socket', ['$rootScope', function($rootScope) {
        var ws;
        var subscribeCallback;
        return {
            connect: function(gadget, errorCallback) {
                if(ws) {
                    ws.close();
                    ws = null;
                }
                ws = new WebSocket("wss://gadgets.dyndns-ip.com/socket?host=" + gadget.host);
                ws.onopen = function() {
                };
                ws.onerror = function() {
                    errorCallback();
                }
                ws.onmessage = function(message) {
                    message = JSON.parse(message.data);
                    var event = message[0];
                    var payload = JSON.parse(message[1]);
                    subscribeCallback(event, payload);
                };
            },
            send: function(message) {
                ws.send(message);
            },
            subscribe: function(callback) {
                subscribeCallback = callback;
            },
            close: function() {
                ws.close();
            }
        }
    }])
    .value('version', '0.1')
    .factory('history', function($rootScope) {
        return {
            getChart: function(series) {
                console.log("series", series);
                config.series = series;
                return config;
            }
        }
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
                    
                });
            }
        }
    });
