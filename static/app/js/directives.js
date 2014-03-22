'use strict';

angular.module('myApp.directives', [])
    .directive('appVersion', ['version', function(version) {
        return function(scope, elm, attrs) {
            elm.text(version);
        };
    }])
    .directive("bootstrapNavbar", ['$location', function($location) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: { gadgets:'=gadgets'},
            templateUrl: "components/navbar.html",
            compile: function(element, attrs) {
                $('[data-hover="dropdown"]').dropdownHover();
            }
        }
    }])
    .directive("methods", ['$modal', 'sockets', function($modal, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: false,
            templateUrl: "components/methods.html",
            controller: function($scope, $timeout, $modal) {
                sockets.subscribe(function (event, message) {
                    if (event == "update" && message.sender == "method runner") {
                        $scope.method = message.method;
                    } else if (event == "method update") {
                        $scope.method.step = message.method.step;
                        $scope.method.time = message.method.time;
                    }
                });
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
                            var name = message.location + " " + message.name;
                        }
                    });
                });

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
                        sockets.send(JSON.stringify(msg));
                    }
                };

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
    }]);


