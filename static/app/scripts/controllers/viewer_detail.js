'use strict';

/**
 * @ngdoc function
 * @name scribbleApp.controller:ViewerDetailCtrl
 * @description
 * # ViewerDetailCtrl
 * Controller of the scribbleApp
 */
angular.module('scribbleApp')
  .controller("ViewerDetailCtrl", ["$scope", "$stateParams", "Restangular",
    function($scope, $stateParams, Restangular) {
      Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
        $scope.currentNote = note;
      });
    }
  ]);
