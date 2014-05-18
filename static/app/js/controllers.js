'use strict';

angular.module('myApp.controllers', [])
    .controller('GadgetsCtrl', ['$rootScope', '$scope', '$routeParams', 'sockets', function($rootScope, $scope, $routeParams, sockets) {
        $scope.locations = {};
        $scope.live = false;
        $scope.name = $routeParams.gadget;
        $scope.host = $routeParams.host;
        sockets.connect($scope.host);
    }])
    .controller('HomeCtrl', ['$rootScope', '$timeout', '$location', function($rootScope, $timeout, $location) {
        
    }])
    .controller('HistoryCtrl', ['$scope', '$http', '$routeParams', 'history', function($scope, $http, $routeParams, history) {
        $scope.gadget = $routeParams.gadget;
    }]);

