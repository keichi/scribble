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
    "ngResource",
    "ngRoute",
    "ngSanitize",
    "ipCookie",
    "ui.router"
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
      });
  })
  .constant("API_ROOT", "http://localhost:8000/api")
;
