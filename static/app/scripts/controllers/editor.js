"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:EditorCtrl
 * @description
 * # EditorCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("EditorCtrl", ["$scope", "$stateParams", "Restangular",
    function($scope, $stateParams, Restangular) {
      Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
        $scope.note = note;
      });

      $scope.save = function() {
        $scope.note.save();
      };

      $scope.aceLoad = function() {
      };

      $scope.aceChange = function() {
      };
    }
  ]);
