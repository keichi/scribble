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
    function ($scope, $stateParams, Restangular) {
      Restangular.all("notes").getList().then(function(notes) {
        $scope.notes = notes;
      });

      Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
        $scope.currentNote = note;
      });
    }
  ]);
