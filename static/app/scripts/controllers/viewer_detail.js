'use strict';

/**
 * @ngdoc function
 * @name scribbleApp.controller:ViewerDetailCtrl
 * @description
 * # ViewerDetailCtrl
 * Controller of the scribbleApp
 */
angular.module('scribbleApp')
  .controller("ViewerDetailCtrl", ["$scope", "$state", "$stateParams", "Restangular",
    function($scope, $state, $stateParams, Restangular) {
      Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
        $scope.currentNote = note;
      });

      $scope.remove = function() {
        $scope.currentNote.remove();
        $state.go("root.viewer");
      };
    }
  ]);
