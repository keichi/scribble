/* global key */
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
            $scope.$parent.removeNote($stateParams.noteId);
            $state.go("root.viewer");
          });
      };

      $scope.$on("viewer.editNote", function() {
        $state.go("root.editor", {noteId: $stateParams.noteId});
      });

      $scope.$on("viewer.deleteNote", function() {
        $scope.remove();
      });

      $scope.$on("viewer.newNote", function() {
        $state.go("root.editor");
      });

      key("up", "viewer", function(e) {
        e.preventDefault();
        $scope.$emit("viewer.selectNote", {
          noteId: $stateParams.noteId,
          direction: "previous"
        });
      });
      key("down", "viewer", function(e) {
        e.preventDefault();
        $scope.$emit("viewer.selectNote", {
          noteId: $stateParams.noteId,
          direction: "next"
        });
      });
    }
  ]);
