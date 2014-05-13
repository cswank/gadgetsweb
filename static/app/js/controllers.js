'use strict';

angular.module('myApp.controllers', []).
    controller('NavbarCtrl', ['$scope', '$modal', 'gadgets', 'auth', function($scope, $modal, gadgets, auth) {

        function login() {
            if ($scope.username == undefined) {
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
            } else {
                auth.login($scope.username, $scope.password, function(){
                    getGadgets();
                });
            }
        }

        function getGadgets() {
            $scope.gadgets = gadgets.get(function(data) {
                $scope.gadgets = data.gadgets;
            }, function() {
                console.log("get gadgets failed");
                login();
            });
        }
        getGadgets();
    }])
    .controller('GadgetsCtrl', ['$rootScope', '$scope', '$routeParams', 'sockets', function($rootScope, $scope, $routeParams, sockets) {
        $scope.locations = {live:false};
        $scope.name = $routeParams.gadget;
        $scope.host = $routeParams.host;
        sockets.connect($scope.host);
    }])
    .controller('HomeCtrl', ['$rootScope', '$timeout', '$location', function($rootScope, $timeout, $location) {
        
    }])
    .controller('HistoryCtrl', ['$scope', '$http', '$routeParams', 'history', function($scope, $http, $routeParams, history) {
        $scope.gadget = $routeParams.gadget;
    }]);

