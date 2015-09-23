/* globals hljs */
"use strict";

/**
 * @ngdoc overview
 * @name scribbleApp
 * @description
 * # scribbleApp
 *
 * Main module of the application.
 */
angular
  .module("scribbleApp", [
    "ngAnimate",
    "ngRoute",
    "ngSanitize",
    "ipCookie",
    "ui.router",
    "restangular",
    "ui.ace",
    "hc.marked",
    "ui.bootstrap",
    "angular-loading-bar",
    "ui-notification"
  ])
  .config(["$stateProvider", "$urlRouterProvider",
    function ($stateProvider, $urlRouterProvider) {
      $urlRouterProvider.otherwise("/");

      $stateProvider
        .state("root", {
          url: "",
          abstract: true,
          views: {
            header: {
              templateUrl: "views/header.html",
              controller: "HeaderCtrl",
              controllerAs: "headerCtrl"
            }
          }
        })
        .state("root.home", {
          url: "/",
          views: {
            "content@": {
              templateUrl: "views/main.html",
              controller: "MainCtrl",
              controllerAs: "mainCtrl"
            }
          }
        })
        .state("root.login", {
          url: "/login",
          views: {
            "content@": {
              templateUrl: "views/login.html",
              controller: "LoginCtrl",
              controllerAs: "loginCtrl"
            }
          }
        })
        .state("root.viewer", {
          url: "/viewer",
          views: {
            "content@": {
              templateUrl: "views/viewer.html",
              controller: "ViewerCtrl",
              controllerAs: "viewerCtrl"
            }
          }
        })
        .state("root.viewer.detail", {
          url: "/:noteId",
          views: {
            detail: {
              templateUrl: "views/viewer_detail.html",
              controller: "ViewerDetailCtrl",
              controllerAs: "viewerDetailCtrl"
            }
          }
        })
        .state("root.editor", {
          url: "/editor/:noteId",
          views: {
            "content@": {
              templateUrl: "views/editor.html",
              controller: "EditorCtrl",
              controllerAs: "editorCtrl"
            }
          }
        });
  }])
  .constant("API_ROOT", "http://localhost:8000/api")
  .config(["RestangularProvider", "API_ROOT",
    function (RestangularProvider, API_ROOT) {
      RestangularProvider.setBaseUrl(API_ROOT);
    }
  ])
  .config(["markedProvider", function(markedProvider) {
    markedProvider.setOptions({
      gfm: true,
      tables: true,
      sanitize: true,
      highlight: function (code, lang) {
        if (!lang || !hljs.getLanguage(lang)) {
          return code;
        }
        return hljs.highlight(lang, code).value;
      }
    });
  }])
  .config(["NotificationProvider", function(NotificationProvider) {
    NotificationProvider.setOptions({
      delay: 5000,
      startTop: 20,
      startRight: 10,
      verticalSpacing: 20,
      horizontalSpacing: 20,
      positionX: 'right',
      positionY: 'bottom'
    });
  }]);
