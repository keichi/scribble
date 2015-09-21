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
      var aceEditor;

      Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
        $scope.note = note;
      });

      var uploadIamge = function(file) {
        $scope.note.one("/images")
          .withHttpConfig({transformRequest: angular.identity})
          .customPOST(file, "", undefined, {"Content-Type": file.type})
          .then(function(resp) {
            aceEditor.insert("![" + resp.uuid + "](" + resp.url + ")");

            // TODO Find better method to force re-render
            $scope.note.content = aceEditor.getValue();
          })
          .finally(function() {
            aceEditor.setReadOnly(false);
          });
      };

      $scope.save = function() {
        $scope.note.save();
      };

      $scope.aceLoaded = function(editor) {
        var event = require("ace/lib/event");
        aceEditor = editor;

        event.addListener(editor.container, "dragover", function(e) {
          var types = e.dataTransfer.types;
          if (types && Array.prototype.indexOf.call(types, 'Files') !== -1) {
            return event.preventDefault(e);
          }
        });

        event.addListener(editor.container, "drop", function(e) {
          var file;
          try {
            file = e.dataTransfer.files[0];
            if (!file.type.startsWith("image/")) {
              throw new Error("Dropped file is not an image");
            }
            if (window.FileReader) {
              editor.setReadOnly(true);
              uploadIamge(file);
            }
            return event.preventDefault(e);
          } catch(err) {
            return event.stopEvent(e);
          }
        });

      };

      $scope.aceChanged = function() {
      };
    }
  ]);
