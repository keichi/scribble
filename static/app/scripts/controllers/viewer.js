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
      $scope.notes = [];
      $scope.currentNoteId = $stateParams.noteId;
      $scope.isBusy = false;

      var nextOffset = 0;
      var nextAvailable = true;
      var pageSize = 10;
      $scope.paginate = function() {
        if (!nextAvailable || $scope.isBusy) {
          return;
        }

        $scope.isBusy = true;
        Restangular.one("my").all("notes")
          .getList({limit: pageSize, offset: nextOffset})
          .then(function(notes) {
            if (notes.length === 0) {
              nextAvailable = false;
            } else {
              Array.prototype.push.apply($scope.notes, notes);
              nextOffset += pageSize;
            }
          })
          .finally(function() {
            $scope.isBusy = false;
          });
      };
    }
  ]);
