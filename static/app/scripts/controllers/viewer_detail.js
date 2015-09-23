"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:ViewerDetailCtrl
 * @description
 * # ViewerDetailCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("ViewerDetailCtrl", ["$scope", "$state", "$stateParams",
    "Restangular", "ModalHelperService",
    function($scope, $state, $stateParams, Restangular, ModalHelperService) {
      Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
        $scope.currentNote = note;
      });

      $scope.remove = function() {
        ModalHelperService.alert("Deleting Note", "Are you sure you want to delete this note?")
          .then(function() {
            $scope.currentNote.remove();
            $state.go("root.viewer");
          });
      };
    }
  ]);
