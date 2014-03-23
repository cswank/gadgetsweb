'use strict';


// Declare app level module which depends on filters, and services
angular.module('Gadgets', [
    'ngRoute',
    'myApp.filters',
    'myApp.services',
    'myApp.directives',
    'myApp.controllers',
    'ui.bootstrap'
]).
config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'partials/home.html', controller: 'HomeCtrl'});
    $routeProvider.when('/gadgets/:gadget', {templateUrl: 'partials/gadgets.html', controller: 'GadgetsCtrl'});
    $routeProvider.when('/history', {templateUrl: 'partials/history.html', controller: 'HistoryCtrl'});
    $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'LoginCtrl'});
    $routeProvider.otherwise({redirectTo: '/'});
}]);
