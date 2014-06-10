'use strict';

angular.module('myApp.directives', [])
    .directive('appVersion', ['version', function(version) {
        return function(scope, elm, attrs) {
            elm.text(version);
        };
    }])
    .directive("bootstrapNavbar", ['$location', 'auth', 'gadgets', 'sockets', function($location, auth, gadgets, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            templateUrl: "components/navbar.html",
            controller: function($scope, $timeout, $modal) {
                
                $('[data-hover="dropdown"]').dropdownHover();
                $scope.login = function() {
                    var dlg = $modal.open({
                        templateUrl: '/dialogs/login.html?c=' + new Date().getTime(),
                        controller: LoginCtrl,
                    });
                    dlg.result.then(function(user) {
                        $scope.username = user.name;
                        $scope.password = user.password;
                        auth.login($scope.username, $scope.password, function(){
                            getGadgets();
                        });
                    });
                }
                
                $scope.logout = function() {
                    auth.logout(function() {
                        $scope.loggedIn = false;
                        sockets.close();
                        $location.url("/");
                        $scope.gadgets = [];
                    });
                }

                $scope.loggedIn = false;

                function getGadgets() {
                    $scope.gadgets = gadgets.get(function(data) {
                        $scope.gadgets = data.gadgets;
                        $scope.loggedIn = true;
                    }, function() {
                        $scope.loggedIn = false;
                        $scope.errMsg = "login failed"
                    });
                }
                getGadgets();
            }
        }
    }])
    .directive("gadgetsConfig", ['$modal', '$http', 'sockets', function($modal, $http, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            templateUrl: "components/gadgets-config.html?x=x",
            controller: function($scope, $timeout, $modal) {
                $scope.cfg = {
                    host:$scope.host,
                    gadgets: []
                };
                $scope.types = [];
                $http.get("api/gadgets/types").success(function(data) {
                    $scope.types = data;
                });
                $scope.newGadget = function() {
                    var dlg = $modal.open({
                        templateUrl: '/dialogs/new-gadget.html?x=y',
                        controller: NewGadgetCtrl,
                        resolve: {
                            types: function () {
                                return $scope.types;
                            }
                        }
                    });
                    
                    dlg.result.then(function(gadget) {
                        $scope.cfg.gadgets.push(gadget);
                    })
                }
                $scope.saveGadgets = function() {
                    $http.post("/api/gadgets", $scope.cfg).success(function(data) {
                        setTimeout(function() {
                            setTimeout(function(){sockets.connect()}, 1200);
                        });
                    });
                }
            }
        }
    }])
    .directive("gadgets", ['$http', '$modal', 'sockets', function($http, $modal, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: {
                locations: "=",
                live: "=",
                name: "="
            },
            templateUrl: "components/gadgets.html",
            link: function($scope, elem, attrs) {
                var promptEvent;
                sockets.subscribe(function (event, message) {
                    if (message.location == "" || message.location == undefined) {
                        return;
                    }
                    $scope.$apply(function() {
                        if (event == "update") {
                            $scope.live = true;
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
    .directive("notes", ['$http', '$modal', 'notes', function($http, $modal, notes) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: {
                name: "="
            },
            templateUrl: "components/notes.html?x=y",
            link: function($scope, elem, attrs) {
                function getNotes() {
                    notes.get($scope.name, function(data) {
                        $scope.notes = data;
                    });
                }
                $scope.addNote = function() {
                    var dlg = $modal.open({
                        templateUrl: '/dialogs/notes.html?c=' + new Date().getTime(),
                        controller: NotesCtrl,
                    });
                    dlg.result.then(function(note) {
                        notes.save($scope.name, note, function() {
                           getNotes(); 
                        });
                    } ,function() {
                        
                    });
                }
                getNotes();
            }
        }
    }])
    .directive("methods", ['$http', '$modal', 'sockets', 'methods', function($http, $modal, sockets, methods) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: false,
            templateUrl: "components/methods.html?x=z",
            controller: function($scope, $timeout, $modal) {
                $scope.recipesAvailable = false;
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
                $http.get("/recipes/_ping").success(function(data) {
                    $scope.recipesAvailable = true;
                });
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
                        templateUrl: '/dialogs/method.html?x=z',
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
                    var e = Math.round(new Date().getTime() / 1000);
                    var s = e - spans[$scope.span];
                    for (var key in $scope.selected.ids) {
                        var val = $scope.selected.ids[key];
                        if (val) {
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
    }])
    .directive('ngEnter', [function() {
        return function(scope, element, attrs) {
            element.bind("keydown keypress", function(event) {
                if(event.which === 13) {
                    scope.$apply(function(){
                        scope.$eval(attrs.ngEnter, {'event': event});
                    });
                    event.preventDefault();
                }
            });
        }
    }])
    .directive('autoFocus', function($timeout) {
        return {
            restrict: 'AC',
            link: function(_scope, _element) {
                $timeout(function(){
                    _element[0].focus();
                }, 10);
            }
        };
    });




        


