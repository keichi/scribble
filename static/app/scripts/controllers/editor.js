/* global key */
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
    "Notification", "ModalHelperService",
    function($scope, $stateParams, Restangular, Notification, modalSvc) {
      var aceEditor = null;
      var isNew = false;

      if ($stateParams.noteId === "") {
        isNew = true;
        $scope.note = {
          title: "",
          content: ""
        };
      } else {
        isNew = false;
        Restangular.one("notes", $stateParams.noteId).get().then(function(note) {
          $scope.note = note;
        });
      }

      var uploadIamge = function(file) {
        if (isNew) {
          modalSvc.alert("Upload Image",
            "Please save this note before uploading an image");
          return;
        }

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
        if ($scope.note.title === "") {
          modalSvc.alert("Failed to save note", "Title is empty.");
          return;
        }

        if (isNew) {
          Restangular.all("notes").post($scope.note)
            .then(function(resp) {
              Notification.success("Note successfully saved.");
              $scope.note = resp;
              isNew = false;
            })
            .catch(function() {
              Notification.error("Note could not be saved.");
            });
        }

        $scope.note.save().then(function() {
          Notification.success("Note successfully saved.");
        })
        .catch(function() {
          Notification.error("Note could not be saved.");
        });
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

        editor.commands.addCommand({
          name: "replace",
          bindKey: {win: "Ctrl-S", mac: "Command-S"},
          exec: function() {
            $scope.$emit("editor.saveNote");
          }
        });
      };

      $scope.aceChanged = function() {
      };

      $scope.$on("editor.saveNote", function() {
        $scope.save();
      });

      key("ctrl+s, command+s", "editor", function(e) {
        e.preventDefault();
        $scope.$broadcast("editor.saveNote");
      });
      key.setScope("editor");
    }
  ]);
