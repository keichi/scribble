"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("MainCtrl", ["$state", "UserService",
    function ($state, userSvc) {
      userSvc.checkLoginStatus().then(function(isLoggedIn) {
        if (isLoggedIn) {
          $state.go("editor");
        }
      });
    }
  ]);
