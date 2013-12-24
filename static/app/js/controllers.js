'use strict';

/* Controllers */

var LoginCtrl = function ($scope, $modalInstance) {
    $scope.user = {
        'name': '',
        'password': ''
    };
    $scope.ok = function () {
        $modalInstance.close($scope.user);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}

var ChartCtrl = function ($scope, $modalInstance, summaries) {
    $scope.summaries = summaries;
    $scope.ok = function () {
        var selected = [];
        for (var i in $scope.summaries) {
            var summary = $scope.summaries[i];
            if (summary.show) {
                selected.push(summary);
            }
        }
        $modalInstance.close(selected);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
    $scope.newValue = function(obj) {
        obj.show = !obj.show;
    };
    
}

var MethodCtrl = function($scope, $modalInstance, method) {
    $scope.method = method;
    var rawMethod = "";
    for (var i in method.method) {
        rawMethod += method.method[i].step + "\n";
    }
    $scope.rawMethod = rawMethod;
    
    $scope.cancel = function(){
        $modalInstance.dismiss('canceled');
    };
    
    $scope.ok = function() {
        var steps = $scope.rawMethod.split("\n");
        var method = [];
        for (var i in steps) {
            method.push({'step': steps[i], 'complete': false})
        }
        $scope.method.method = method;
        $modalInstance.close($scope.method);
    };
};

angular.module('myApp.controllers', []).
    controller('GadgetsCtrl', ['$scope', '$http', '$timeout', '$modal', '$location', 'socket', function($scope, $http, $timeout, $modal, $location, socket) {
        $scope.gadget = {'name': 'select a host', 'host': 'remove me'}
        $scope.method = {'name': 'select a method', 'steps': []}
        var events = {};
        var promptEvent;
        $http.get('/methods').success(function (data, status, headers, config) {
            $scope.methods = [$scope.method];
            for (var i in data.methods) {
                var rawMethod = data.methods[i];
                var steps = [];
                for (var j in rawMethod.steps) {
                    steps.push({id: rawMethod.id, step:rawMethod.steps[j], complete:false})
                }
                $scope.methods.push({id: rawMethod.id, name: rawMethod.name, steps: steps});
            }
            console.log($scope.methods);
        });

        $http.get('/gadgets').success(function (data, status, headers, config) {
            data.gadgets.unshift($scope.gadget);
            $scope.gadgets = data.gadgets;
        });

        $scope.runMethod = function() {
            var method = [];
            for (var i in $scope.method.steps) {
                method.push($scope.method.steps[i].step);
            }
            var msg = {event: 'method', 'message': {method:method}};
            socket.send(JSON.stringify(msg));
        };

        var saveMethod = function() {
            var method = [];
            for (var i in $scope.method.method) {
                method.push($scope.method.steps[i].step);
            }
            var url, httpMethod, data
            if ($scope.method.id != undefined && $scope.method.id > 0) {
                url = '/methods/' + $scope.method.id.toString();
                httpMethod = 'PUT'
                data = {id: $scope.method.id, name: $scope.method.name, steps:method};
            } else {
                url = '/methods';
                httpMethod = 'POST'
                data = {name: $scope.method.name, steps:method};
            }
            $http({
                url: url,
                method: httpMethod,
                data: JSON.stringify(data),
                headers: {'Content-Type': 'application/json'}
            }).success(function (data, status, headers, config) {
                console.log(data, status);
            }).error(function (data, status, headers, config) {
                
            });
        };

        $scope.addMethod = function() {
            var dlg = $modal.open({
                templateUrl: '/dialogs/method.html?x=1',
                controller: MethodCtrl,
                resolve: {
                    method: function () {
                        return $scope.method;
                    }
                }
            });
            dlg.result.then(function(method) {
                saveMethod();
                $scope.method = method;
            } ,function(){
                
            });
        };
        
        $scope.logout = function() {
            $http({
                url: '/logout',
                method: "POST",
                headers: {'Content-Type': 'application/json'}
            }).success(function (data, status, headers, config) {
                socket.close();
                $scope.locations = {};
                getCredentials();
            }).error(function (data, status, headers, config) {
                
            });
        }

        $scope.clearDummyMethod = function() {
            console.log($scope.methods);
            if ($scope.methods[0].name == 'select a method') {
                $scope.methods.shift();
            }
            console.log($scope.methods);
        };

        $scope.connect = function() {
            if ($scope.gadgets[0].host == 'remove me') {
                $scope.gadgets.shift();
            }
            socket.connect($scope.gadget, getCredentials);
        };
        var doLogin = function() {
            $http({
                url: '/login',
                method: "POST",
                data: JSON.stringify({username:$scope.username, password: $scope.password}),
                headers: {'Content-Type': 'application/json'}
            }).success(function (data, status, headers, config) {
                socket.connect($scope.gadget, getCredentials);
            }).error(function (data, status, headers, config) {
            });
        };
        
        var getCredentials = function() {
            var dlg = $modal.open({
                templateUrl: '/dialogs/login.html?c=' + new Date().getTime(),
                controller: LoginCtrl,
            });
            dlg.result.then(function(user) {
                $scope.username = user.name;
                $scope.password = user.password;
                doLogin();
            } ,function(){
                
            });
        };
        
        $scope.login = function() {
            $scope.loginPromptShouldBeOpen = false;
            doLogin();
        }

        function getCommandValue(value) {
            var commandValue;
            if (value) {
                commandValue = 'off';
            }
            else {
                commandValue = 'on';
            }
            return commandValue;
        }

        socket.subscribe(function (event, message) {
            $scope.$apply(function() {
                if (event == "UPDATE" || event == "status") {
                    $scope.locations = message.locations;
                    for (var locationKey in message.locations) {
                        var location = message.locations[locationKey];
                        for (var deviceKey in location.output) {
                            var device = location.output[deviceKey];
                            device.key = deviceKey;
                        }
                    }
                    var method = message.method;
                    if (method != undefined && method.method != undefined && method.method.length > 0) {
                        var step, countdown;
                        var steps = [];
                        for (var i=0; i < method.method.length; i++) {
                            step = {step: method.method[i]};
                            step.complete = (i < method.step) ? 'step-complete' : 'step-incomplete';
                            steps.push(step);
                        }
                        countdown = (method.countdown > 0) ? 'countdown: ' + method.countdown.toString() : ''
                        $scope.method.step =  method.step;
                        $scope.method.steps =  steps;
                        $scope.method.countdown = countdown;
                    }
                }
                else if (event == "commands") {
                    $scope.commands = message;
                }
            });
        });

        $scope.sendCommand = function() {
            $scope.promptShouldBeOpen = false;
            var command = $scope.currentCommand + $scope.commandArgument;
            var msg = {event: command, 'message': {}};
            socket.send(JSON.stringify(msg));
            $scope.currentCommand = null;
            $scope.commandArgument = null;
        };

        $scope.getArguments = function(location, device, value) {
            promptEvent = $timeout(function() {
                var commandValue = getCommandValue(value);
                $scope.currentCommand = events[location][device][commandValue];
                $scope.promptShouldBeOpen = true;
            }, 1000);
        };

        $scope.toggle = function(location, device, value) {
            
            $timeout.cancel(promptEvent);
            if (!$scope.promptShouldBeOpen) {
                var commandValue = getCommandValue(value);
                var command = $scope.commands[location][device][commandValue];
                var msg = {event:command, message:{}};
                socket.send(JSON.stringify(msg));
            }
        };
    }]).
    controller('HistoryCtrl', ['$scope', '$http', '$modal', 'history', function($scope, $http, $modal, history) {
        $scope.promptShouldBeOpen = false;
        $scope.openPrompt = function(val) {
            $scope.promptShouldBeOpen = val;
        }
        $http.get("/history/locations/summary").success( function(data) {
            var summary;
            for (var i in data) {
                summary = data[i];
                summary.selected = false;
                summary.show = false;
            }
            $scope.summaries = data;
        });

        $scope.choose = function() {
            var dlg = $modal.open({
                templateUrl: '/dialogs/chart.html',
                controller: ChartCtrl,
                resolve: {
                    summaries: function () {
                        return $scope.summaries;
                    }
                }
            });
            dlg.result.then(function(selected) {
                var now = new Date().getTime();
                var start = now - 604800000; //one week
                //var start = now - (86400000 * 2); //two days
                var query = {start: start, end: now}
                var url;
                var summary;
                var chartData = [];
                
                for (var i in selected) {
                    summary = selected[i];
                    url = '/history/locations/' + summary.location + '/directions/'  + summary.direction + '/devices/' + summary.name;
                    $http({method:'GET', url:url, params:query}).success(function(data) {
                        chartData.push(data);
                        if (chartData.length == selected.length) {
                            $scope.history = history.getChart(chartData);
                        }
                    })
                }
            } ,function(){
                
            });
        }
    }])
