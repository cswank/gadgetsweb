'use strict';

angular.module('myApp.controllers', []).
    controller('NavbarCtrl', ['$rootScope', '$scope', 'gadgets', 'auth', function($rootScope, $scope, gadgets, auth) {
        
        $rootScope.$on("login", function(event){
            login();
        });
        
        $scope.gadgets = gadgets.get(function(data) {
            console.log(data);
            $scope.gadgets = data.gadgets;
        });
        
        function login() {
            if ($scope.username == undefined) {
                var dlg = $modal.open({
                    templateUrl: '/dialogs/login.html?c=' + new Date().getTime(),
                    controller: LoginCtrl,
                });
                dlg.result.then(function(user) {
                    $scope.username = user.name;
                    $scope.password = user.password;
                    auth.login($scope.username, $scope.password);
                } ,function(){});
            } else {
                auth.login($scope.username, $scope.password);
            }
        }
        
    }]).
    controller('GadgetsCtrl', ['$rootScope', '$scope', '$routeParams', 'sockets', function($rootScope, $scope, $routeParams, sockets) {
        $scope.name = $routeParams.gadget;
        $scope.host = $routeParams.host;
        sockets.connect($scope.host);
    }])
    .controller('HistoryCtrl', [function() {
        
    }]);
