'use strict';

angular.module('myApp.directives', [])
    .directive('appVersion', ['version', function(version) {
        return function(scope, elm, attrs) {
            elm.text(version);
        };
    }])
    .directive("bootstrapNavbar", ['$location', 'auth', 'sockets', function($location, auth, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: { gadgets:'=gadgets'},
            templateUrl: "components/navbar.html",
            controller: function($scope, $timeout, $modal) {
                $('[data-hover="dropdown"]').dropdownHover();
                $scope.logout = function() {
                    auth.logout(function(){
                        sockets.close();
                        $location.url("/");
                    });
                }
            }
        }
    }])
    .directive("gadgetsConfig", ['$modal', '$http', function($modal, $http) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: false,
            templateUrl: "components/gadgets-config.html",
            controller: function($scope, $timeout, $modal) {
                $scope.types = [];
                $http.get("api/gadgets/types").success(function(data) {
                    $scope.types = data;
                });
                $scope.newGadget = function() {
                    var dlg = $modal.open({
                        templateUrl: '/dialogs/new-gadget.html?x=x',
                        controller: NewGadgetCtrl,
                        resolve: {
                            types: function () {
                                return $scope.types;
                            }
                        }
                    });
                    
                    dlg.result.then(function(gadget) {
                        console.log(gadget);
                    })
                }
            }
        }
    }])
    .directive("gadgets", ['$modal', 'sockets', function($modal, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: {
                locations: "="
            },
            templateUrl: "components/gadgets.html?x=x",
            link: function($scope, elem, attrs) {
                var promptEvent;
                sockets.subscribe(function (event, message) {
                    if (message.location == "" || message.location == undefined) {
                        return;
                    }
                    $scope.$apply(function() {
                        if (event == "update") {
                            $scope.locations.live = true;
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
                        }
                    });
                });
                $scope.sendCommand = function() {
                    $scope.promptShouldBeOpen = false;
                    var command = $scope.currentCommand + $scope.commandArgument;
                    var msg = {event: command, 'message': {}};
                    sockets.send(JSON.stringify(msg));
                    $scope.currentCommand = null;
                    $scope.commandArgument = null;
                };
            }
        }
    }])
    .directive("methods", ['$http', '$modal', 'sockets', 'methods', function($http, $modal, sockets, methods) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: false,
            templateUrl: "components/methods.html?x=x",
            controller: function($scope, $timeout, $modal) {
                function getMethods() {
                    $scope.method = {id:-1,name:"select"};
                    $scope.methods = [$scope.method];
                    methods.get($scope.name, function(data) {
                        $scope.showMethods = true;
                        for (var i in data.methods) {
                            var rawMethod = data.methods[i];
                            $scope.methods.push(rawMethod);
                        }
                    })
                }
                getMethods();
                sockets.subscribe(function (event, message) {
                    if (event == "update" && message.sender == "method runner") {
                        if (message.method.steps == null) {
                            $scope.method = {id:-1,name:"select"};
                            $scope.methods[0] = $scope.method;
                        } else {
                            $scope.method = message.method;
                        }
                    } else if (event == "method update") {
                        $scope.$apply(function() {
                            $scope.method.step = message.method.step;
                            $scope.method.time = message.method.time;
                        });
                    }
                });

                $scope.runMethod = function() {
                    var msg = {event: 'method', message: {type: 'method', method:$scope.method}};
                    sockets.send(JSON.stringify(msg));
                };

                $scope.confirm = function(step) {
                    var msg = {
                        event: 'update',
                        message: {
                            type: 'update',
                            body:step,
                        }
                    };
                    sockets.send(JSON.stringify(msg));
                };
                
                $scope.checkUserPrompt = function(i) {
                    var step = $scope.method.steps[i];
                    return step != undefined && step.indexOf("wait for user") == 0 && i == $scope.method.step;
                };
                
                $scope.addMethod = function() {
                    if ($scope.method.name == "select") {
                        $scope.method = {name:""};
                    }
                    var dlg = $modal.open({
                        templateUrl: '/dialogs/method.html?x=x',
                        controller: MethodCtrl,
                        resolve: {
                            method: function () {
                                return $scope.method;
                            }
                        }
                    });
                    dlg.result.then(function(method) {
                        methods.save($scope.name, method);
                        $scope.method = method;
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

                $scope.clearMethod = function() {
                    var msg = {event: 'command', message: {type: 'command', body:'clear method'}};
                    sockets.send(JSON.stringify(msg));
                    getMethods();
                };

                $scope.deleteMethod = function() {
                    methods.delete($scope.name, $scope.method, function() {
                        getMethods();
                    });
                }
            }
        }
    }])
    .directive('toggleSwitch', ['sockets', '$modal', '$timeout', function (sockets, $modal, $timeout) {
        return {
            restrict: 'EA',
            replace: true,
            scope: {
                device: '=',
                disabled: '@',
                onLabel: '@',
                offLabel: '@',
                knobLabel: '@'
            },
            template: '<div class="switch" ng-mousedown="getArguments()" ng-click="toggle()" ng-class="{ \'disabled\': disabled }"><div class="switch-animate" ng-class="{\'switch-off\': !device.value.value, \'switch-on\': device.value.value}"><span class="switch-left" ng-bind="onLabel"></span><span class="knob" ng-bind="knobLabel"></span><span class="switch-right" ng-bind="offLabel"></span></div></div>',
            controller: function($scope) {
                var promptEvent;
                $scope.toggle = function() {
                    $timeout.cancel(promptEvent);
                    var command;
                    if (!$scope.device.value.value) {
                        command = $scope.device.info.on;
                    } else {
                        command = $scope.device.info.off;
                    }
                    var msg = {event:'command', message:{type:'command', body:command}};
                    sockets.send(JSON.stringify(msg));
                };
                $scope.getArguments = function() {
                    promptEvent = $timeout(function() {
                        var dlg = $modal.open({
                            templateUrl: '/dialogs/command.html?c=' + new Date().getTime(),
                            controller: CommandCtrl,
                            resolve: {
                                command: function () {
                                    return $scope.device.info.on;
                                }
                            }
                        });
                        dlg.result.then(function(command) {
                            var msg = {event:'command', message:{type:'command', body:command}};
                            sockets.send(JSON.stringify(msg));
                        } ,function(){
                            
                        });
                    }, 1000);
                };
            },
            compile: function(element, attrs) {
                if (!attrs.onLabel) { attrs.onLabel = 'On'; }
                if (!attrs.offLabel) { attrs.offLabel = 'Off'; }
                if (!attrs.knobLabel) { attrs.knobLabel = '\u00a0'; }
                if (!attrs.disabled) { attrs.disabled = false; }
            },
        };
    }])
    .directive('historyChart', ['$http', 'history', function ($http, history) {
        var spans = {
            "hour": 60 * 60,
            "day": 24 * 60 * 60,
            "week": 7 * 24 * 60 * 60,
        }
        return {
            restrict: 'E',
            replace: true,
            scope: {
                gadget: '=gadget',
            },
            templateUrl: "components/history.html?x=x",
            controller: function($scope) {
                $scope.span = "hour";
                Highcharts.setOptions({
	            global: {
		        useUTC: false
	            }
                });
                $scope.chartConfig = {
                    options: {
                        chart: {
                            type: 'line',
                            zoomType: 'x'
                        }
                    },
                    plotOptions: {
                        line: {
                            marker: {
                                enabled: false
                            }
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
                $scope.selected = {
                    ids:{}
                };
                history.getDevices($scope.gadget, function(data){
                    $scope.links = data.links;
                });
                $scope.getHistory = function(){
                    var series = [];
                    
                    for (var key in $scope.selected.ids) {
                        var val = $scope.selected.ids[key];
                        if (val) {
                            var e = Math.round(new Date().getTime() / 1000);
                            var s = e - spans[$scope.span];
                            var url = key + '?start=' + s + '&end=' + e;
                            $http.get(url).success(function(data) {
                                series.push(data);
                            });
                        }
                    }
                    $scope.chartConfig.series = series;
                }
            }
        }
    }]);



        


