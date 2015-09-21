"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:ViewerCtrl
 * @description
 * # ViewerCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("ViewerCtrl", ["$scope", "$stateParams", "Restangular",
    function ($scope, $stateParams, Restangular) {
      Restangular.one("my").all("notes").getList().then(function(notes) {
        $scope.notes = notes;
      });

      $scope.currentNoteId = $stateParams.noteId;
    }
  ]);
