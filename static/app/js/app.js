'use strict';


// Declare app level module which depends on filters, and services
angular.module('Gadgets', [
    'ngRoute',
    'myApp.filters',
    'myApp.services',
    'myApp.directives',
    'myApp.controllers',
    'ui.bootstrap',
    'highcharts-ng',
    'ng-context-menu',
    'ngStorage'
]).
config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'partials/home.html', controller: 'HomeCtrl'});
    $routeProvider.when(
        '/gadgets/:gadget',
        {
            templateUrl: 'partials/gadgets.html',
            controller: 'GadgetsCtrl'
        }
    );
    $routeProvider.when('/history/:gadget', {templateUrl: 'partials/history.html', controller: 'HistoryCtrl'});
    $routeProvider.otherwise({redirectTo: '/'});
}]);
