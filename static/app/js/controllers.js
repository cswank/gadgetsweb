
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

var CommandCtrl = function ($scope, $modalInstance, command) {
    $scope.command = {
        command: command,
        arg: "",
    };
    $scope.ok = function () {
        var cmd = $scope.command.command + " " + $scope.command.arg;
        $modalInstance.close(cmd);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}

var RecipeCtrl = function ($scope, $modalInstance) {
    $scope.recipe = {
        name: "",
        grainTemperature: "",
    };
    $scope.ok = function () {
        $modalInstance.close($scope.recipe);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}

var ChartCtrl = function ($scope, $modalInstance, links) {
    $scope.links = links;
    $scope.ok = function() {
        var selected = [];
        for (var i in $scope.links) {
            var link = $scope.links[i];
            if (link.selected) {
                selected.push(link);
            }
        }
        $modalInstance.close(selected);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
    $scope.newValue = function(obj) {
        obj.selected = !obj.selected;
    };
}

var MethodCtrl = function($scope, $modalInstance, method) {
    $scope.method = method;
    var rawMethod = "";
    for (var i in method.steps) {
        rawMethod += method.steps[i] + "\n";
    }
    $scope.rawMethod = rawMethod;
    
    $scope.cancel = function(){
        $modalInstance.dismiss('canceled');
    };
    
    $scope.ok = function() {
        var steps = $scope.rawMethod.split("\n");
        $scope.method.steps = steps;
        $modalInstance.close($scope.method);
    };
};

angular.module('myApp.controllers', []).
    controller('GadgetsCtrl', ['$scope', '$http', '$timeout', '$modal', '$location', 'socket', function($scope, $http, $timeout, $modal, $location, socket) {
        $scope.showMethods = false;
        $scope.gadget = {'name': 'select a host', 'host': 'remove me'};
        $scope.method = {'name': 'select a method', 'steps': []};
        var events = {};
        var promptEvent;
        

        $http.get('/gadgets').success(function (data, status, headers, config) {
            data.gadgets.unshift($scope.gadget);
            $scope.gadgets = data.gadgets;
        });

        $scope.runMethod = function() {
            var msg = {event: 'method', message: {type: 'method', method:$scope.method}};
            socket.send(JSON.stringify(msg));
        };

        $scope.clearMethod = function() {
            var msg = {event: 'command', message: {type: 'command', body:'clear method'}};
            socket.send(JSON.stringify(msg));
        };

        var saveMethod = function() {
            var url, httpMethod, data
            if ($scope.method.id != undefined && $scope.method.id > 0) {
                url = '/gadgets/' + $scope.gadget.name + '/methods/' + $scope.method.id.toString();
                httpMethod = 'PUT'
            } else {
                url = '/gadgets/' + $scope.gadget.name + '/methods';
                httpMethod = 'POST'
            }
            $http({
                url: url,
                method: httpMethod,
                data: JSON.stringify($scope.method),
                headers: {'Content-Type': 'application/json'}
            }).success(function (data, status, headers, config) {
                
            }).error(function (data, status, headers, config) {
                console.log("error saving method");
            });
        };

        $scope.addMethod = function() {
            var dlg = $modal.open({
                templateUrl: '/dialogs/method.html',
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

        $scope.getRecipe = function() {
            var dlg = $modal.open({
                templateUrl: '/dialogs/recipe.html?c=' + new Date().getTime(),
                controller: RecipeCtrl,
            });
            dlg.result.then(function(recipe) {
                var url = '/recipes/' + recipe.name + '?grainTemperature=' + recipe.grainTemperature;
                $http.get(url).success(function (data, status, headers, config) {
                    $scope.method = data;
                });
            } ,function() {
                
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
            if ($scope.methods[0].name == 'select a method') {
                $scope.methods.shift();
            }
        };

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

        $scope.connect = function() {
            getMethods();
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
        $scope.locations = {};
        socket.subscribe(function (event, message) {
            $scope.$apply(function() {
                if (event == "update" && message.sender == "method runner") {
                    $scope.method = message.method;
                } else if (event == "update") {
                    if ($scope.locations[message.location] == undefined) {
                        $scope.locations[message.location] = {};
                    }
                    if ($scope.locations[message.location][message.name] == undefined) {
                        $scope.locations[message.location][message.name] = {};
                    }
                    if ($scope.locations[message.location][message.name]['value'] != undefined) {
                        $scope.locations[message.location][message.name]['value'] = message.value;
                    } else {
                        $scope.locations[message.location][message.name] = message;
                    }
                } else if (event == "method update") {
                    $scope.method.step = message.method.step;
                    $scope.method.time = message.method.time;
                }
            });
        });

        $scope.checkUserPrompt = function(i) {
            var step = $scope.method.steps[i];
            return step != undefined && step.indexOf("wait for user") == 0 && i == $scope.method.step;
        };

        $scope.confirm = function(step) {
            var msg = {
                event: 'update',
                message: {
                    type: 'update',
                    body:step,
                }
            };
            socket.send(JSON.stringify(msg));
            
        }
        
        $scope.sendCommand = function() {
            $scope.promptShouldBeOpen = false;
            var command = $scope.currentCommand + $scope.commandArgument;
            var msg = {event: command, 'message': {}};
            socket.send(JSON.stringify(msg));
            $scope.currentCommand = null;
            $scope.commandArgument = null;
        };

        $scope.getArguments = function(device) {
            console.log(device);
            promptEvent = $timeout(function() {
                var dlg = $modal.open({
                    templateUrl: '/dialogs/command.html?c=' + new Date().getTime(),
                    controller: CommandCtrl,
                    resolve: {
                        command: function () {
                            return device.info.on;
                        }
                    }
                });
                dlg.result.then(function(command) {
                    console.log(command);
                    var msg = {event:'command', message:{type:'command', body:command}};
                    socket.send(JSON.stringify(msg));
                } ,function(){
                    
                });
            }, 1000);
        };

        $scope.toggle = function(device) {
            $timeout.cancel(promptEvent);
            if (!$scope.promptShouldBeOpen) {
                var command;
                if (!device.value.value) {
                    command = device.info.on;
                } else {
                    command = device.info.off;
                }
                var msg = {event:'command', message:{type:'command', body:command}};
                socket.send(JSON.stringify(msg));
            }
        };
    }]).
    controller('HistoryCtrl', ['$scope', '$http', '$modal', 'history', function($scope, $http, $modal, history) {
        $scope.promptShouldBeOpen = false;
        $scope.openPrompt = function(val) {
            $scope.promptShouldBeOpen = val;
        }
        $http.get("/history/devices").success( function(data) {
            for (var i in data.links) {
                var link = data.links[i];
                link.selected = false;
            }
            $scope.links = data.links;
        });

        $scope.chartConfig = {
            options: {
                chart: {
                    type: 'line',
                    zoomType: 'x'
                }
            },
            series: [],
            title: {
                text: 'Gadgets'
            },
            xAxis: {
                type: 'datetime',
                dateTimeLabelFormats: { // don't display the dummy year
                    month: '%e. %b',
                    year: '%b'
                }
            },
            loading: false
        }
        
        $scope.choose = function() {
            var dlg = $modal.open({
                templateUrl: '/dialogs/chart.html?x=yyy',
                controller: ChartCtrl,
                resolve: {
                    links: function () {
                        return $scope.links;
                    }
                }
            });
            dlg.result.then(function(selected) {
                $scope.chartConfig.series = [];
                var now = new Date().getTime();
                var start = now - 604800000; //one week
                //var start = now - (86400000 * 2); //two days
                var query = {
                    start: Math.round(start / 1000),
                    end: Math.round(now / 1000)
                }
                var link;
                for (var i in selected) {
                    link = selected[i];
                    $http({method:'GET', url:link.path, params:query}).success(function(data) {
                        $scope.chartConfig.series.push(data[0]);
                    });
                }
            } ,function(){
                
            });
        }
    }])
