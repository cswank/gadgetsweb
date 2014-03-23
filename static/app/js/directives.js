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
                    console.log("loggin out");
                    auth.logout(function(){
                        console.log("logged out");
                        sockets.close();
                        $location.url("/#/");
                    });
                }
            }
        }
    }])
    .directive("methods", ['$http', '$modal', 'sockets', 'methods', function($http, $modal, sockets, methods) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: false,
            templateUrl: "components/methods.html",
            controller: function($scope, $timeout, $modal) {
                methods.get($scope.name, function(data) {
                    $scope.showMethods = true;
                    $scope.methods = [$scope.method];
                    for (var i in data.methods) {
                        var rawMethod = data.methods[i];
                        $scope.methods.push(rawMethod);
                    }
                })
                
                sockets.subscribe(function (event, message) {
                    if (event == "update" && message.sender == "method runner") {
                        $scope.method = message.method;
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
                        methods.save($scope.name, method);
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

                $scope.clearMethod = function() {
                    var msg = {event: 'command', message: {type: 'command', body:'clear method'}};
                    console.log(msg);
                    sockets.send(JSON.stringify(msg));
                };
            }
        }
    }])
    .directive("gadgets", ['$modal', 'sockets', function($modal, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: false,
            templateUrl: "components/gadgets.html?x=x",
            controller: function($scope, $timeout, $modal) {
                var promptEvent;
                $scope.locations = {};
                sockets.subscribe(function (event, message) {
                    if (message.location == "") {
                        return;
                    }
                    $scope.$apply(function() {
                        if (event == "update") {
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

                $scope.getArguments = function(device) {
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
                            var msg = {event:'command', message:{type:'command', body:command}};
                            sockets.send(JSON.stringify(msg));
                        } ,function(){
                            
                        });
                    }, 1000);
                };
            }
        }
    }])
    .directive('toggleSwitch', ['sockets', function (sockets) {
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
            template: '<div class="switch" ng-click="toggle()" ng-class="{ \'disabled\': disabled }"><div class="switch-animate" ng-class="{\'switch-off\': !device.value.value, \'switch-on\': device.value.value}"><span class="switch-left" ng-bind="onLabel"></span><span class="knob" ng-bind="knobLabel"></span><span class="switch-right" ng-bind="offLabel"></span></div></div>',
            controller: function($scope) {
                // $scope.toggle = function(device) {
                //     $timeout.cancel(promptEvent);
                //     if (!$scope.promptShouldBeOpen) {
                //         var command;
                //         if (!device.value.value) {
                //             command = device.info.on;
                //         } else {
                //             command = device.info.off;
                //         }
                //         var msg = {event:'command', message:{type:'command', body:command}};
                //         sockets.send(JSON.stringify(msg));
                //     }
                // };
                $scope.toggle = function() {
                    var command;
                    if (!$scope.device.value.value) {
                        command = $scope.device.info.on;
                    } else {
                        command = $scope.device.info.off;
                    }
                    var msg = {event:'command', message:{type:'command', body:command}};
                    sockets.send(JSON.stringify(msg));
                };
            },
            compile: function(element, attrs) {
                if (!attrs.onLabel) { attrs.onLabel = 'On'; }
                if (!attrs.offLabel) { attrs.offLabel = 'Off'; }
                if (!attrs.knobLabel) { attrs.knobLabel = '\u00a0'; }
                if (!attrs.disabled) { attrs.disabled = false; }
            },
        };
    }]);
        


