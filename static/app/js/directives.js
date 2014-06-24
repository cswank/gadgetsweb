'use strict';

angular.module('myApp.directives', [])
    .directive('appVersion', ['version', function(version) {
        return function(scope, elm, attrs) {
            elm.text(version);
        };
    }])
    .directive("bootstrapNavbar", ['$location', '$localStorage', 'auth', 'gadgets', 'sockets', function($location, $localStorage, auth, gadgets, sockets) {
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            templateUrl: "components/navbar.html?x=z",
            controller: function($scope, $timeout, $modal) {
                $scope.$storage = $localStorage.$default({
                    username: ""
                });
                $scope.username = "";
                $('[data-hover="dropdown"]').dropdownHover();
                $scope.login = function(errorMessage) {
                    var dlg = $modal.open({
                        templateUrl: '/dialogs/login.html?c=' + new Date().getTime(),
                        controller: LoginCtrl,
                        resolve: {
                            message: function () {
                                return errorMessage;
                            }
                        }
                    });
                    dlg.result.then(function(user) {
                        $scope.username = user.name;
                        $scope.$storage.username = user.name;
                        $scope.password = user.password;
                        auth.login($scope.username, $scope.password, function(){
                            getGadgets();
                        }, function(){
                            $scope.login("username or password not correct, please try again");
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
        function doUpdate($scope, event, message) {
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
        }
        return {
            restrict: "E",
            replace: true,
            transclude: true,
            scope: {
                locations: "=",
                live: "=",
                name: "=",
                mobile: "="
            },
            templateUrl: "components/gadgets.html",
            link: function($scope, elem, attrs) {
                var promptEvent;
                sockets.subscribe(function (event, message) {
                    if (message.location == "" || message.location == undefined) {
                        return;
                    }
                    if ($scope.mobile) {
                        doUpdate($scope, event, message);
                    } else {
                        $scope.$apply(function() {
                            doUpdate($scope, event, message);
                        });
                    }
                    
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
            templateUrl: 'components/toggle-switch.html',
            controller: function($scope) {
                $scope.mobile = (function() {
                    var check = false;
                    (function(a){if(/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(a)||/1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(a.substr(0,4)))check = true})(navigator.userAgent||navigator.vendor||window.opera);
                    return check;
                })();
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




        


