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
    "restangular"
  ])
  .config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $stateProvider
      .state("home", {
        url: "/",
        templateUrl: "views/main.html",
        controller: "MainCtrl",
        controllerAs: "mainCtrl"
      })
      .state("login", {
        url: "/login",
        templateUrl: "views/login.html",
        controller: "LoginCtrl",
        controllerAs: "loginCtrl"
      })
      .state("editor", {
        url: "/editor/:noteId",
        templateUrl: "views/editor.html",
        controller: "EditorCtrl",
        controllerAs: "editorCtrl"
      });
  })
  .constant("API_ROOT", "http://localhost:8000/api")
  .config(["RestangularProvider", "API_ROOT",
    function (RestangularProvider, API_ROOT) {
      RestangularProvider.setBaseUrl(API_ROOT);
      RestangularProvider.setRequestSuffix("");
    }
  ]);
