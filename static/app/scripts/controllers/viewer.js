/* global _, key */
"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:ViewerCtrl
 * @description
 * # ViewerCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("ViewerCtrl", ["$scope", "$state", "Restangular",
    function ($scope, $state, Restangular) {
      $scope.notes = [];
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

      $scope.removeNote = function(noteId) {
        _.remove($scope.notes, function(note) {
          return note.id === parseInt(noteId, 10);
        });
      };

      $scope.$on("viewer.selectNote", function(e, data) {
        var noteId = parseInt(data.noteId, 10);
        var direction = data.direction;

        var currentIdx = 0;
        if ($state.is("root.viewer.detail")) {
          currentIdx = _.findIndex($scope.notes, function(note) {
            return note.id === noteId;
          });

          if (direction === "next") {
            currentIdx++;
          } else if (direction === "previous") {
            currentIdx--;
          }
        }
        if (0 <= currentIdx && currentIdx < $scope.notes.length) {
          $state.go("root.viewer.detail", {noteId: $scope.notes[currentIdx].id});
        }
      });

      key("enter", "viewer", function(e) {
        e.preventDefault();
        $scope.$broadcast("viewer.editNote");
      });
      key("backspace", "viewer", function(e) {
        e.preventDefault();
        $scope.$broadcast("viewer.deleteNote");
      });
      key.setScope("viewer");
    }
  ]);
