'use strict';

var LoginCtrl = function ($scope, $modalInstance) {
    $scope.user = {
        'name': '',
        'password': ''
    };
    $scope.ok = function () {
        $modalInstance.close($scope.user);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}

var CommandCtrl = function ($scope, $modalInstance, command) {
    $scope.command = {
        command: command,
        arg: "",
    };
    $scope.ok = function () {
        var cmd = $scope.command.command + " " + $scope.command.arg;
        $modalInstance.close(cmd);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}

var RecipeCtrl = function ($scope, $modalInstance) {
    $scope.recipe = {
        name: "",
        grainTemperature: "",
    };
    $scope.ok = function () {
        $modalInstance.close($scope.recipe);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}

var NewGadgetCtrl = function ($scope, $modalInstance, types) {
    $scope.types = types;
    $scope.type = {};
    $scope.gadget = {
        name: "",
        location: "",
    };
    $scope.ok = function () {
        $modalInstance.close($scope.gadget);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
    $scope.select = function() {
        $scope.selectedType = $scope.types[$scope.gadget.type];
        console.log($scope.selectedType);
    }
}

var ChartCtrl = function ($scope, $modalInstance, links) {
    $scope.links = links;
    $scope.ok = function() {
        var selected = [];
        for (var i in $scope.links) {
            var link = $scope.links[i];
            if (link.selected) {
                selected.push(link);
            }
        }
        $modalInstance.close(selected);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
    $scope.newValue = function(obj) {
        obj.selected = !obj.selected;
    };
}

var MethodCtrl = function($scope, $modalInstance, method) {
    $scope.method = method;
    var rawMethod = "";
    for (var i in method.steps) {
        rawMethod += method.steps[i] + "\n";
    }
    $scope.rawMethod = rawMethod;
    
    $scope.cancel = function(){
        $modalInstance.dismiss('canceled');
        $scope.method = {id:-1,name:"select"};
    };
    
    $scope.ok = function() {
        var steps = $scope.rawMethod.split("\n");
        $scope.method.steps = steps;
        $modalInstance.close($scope.method);
    };
};
