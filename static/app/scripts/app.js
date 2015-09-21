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
    "ui.bootstrap"
  ])
  .config(["$stateProvider", "$urlRouterProvider",
    function ($stateProvider, $urlRouterProvider) {
      $urlRouterProvider.otherwise("/");

      $stateProvider
        .state("root", {
          url: "",
          views: {
            header: {
              templateUrl: "views/header.html"
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
      highlight: function (code) {
        return hljs.highlightAuto(code).value;
      }
    });
  }]);
