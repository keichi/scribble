/* globals md5 */
"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:HeaderCtrl
 * @description
 * # HeaderCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("HeaderCtrl", ["$scope", "$state", "UserService",
    function ($scope, $state, userSvc) {
      userSvc.checkLoginStatus();

      $scope.$watch(function() {
        return userSvc.user;
      }, function() {
        if (userSvc.user) {
          $scope.isLoggedIn = true;
          $scope.user = userSvc.user;

          var hash = md5(userSvc.user.email.trim().toLowerCase());
          $scope.gravatarHash = hash;
        } else {
          $scope.isLoggedIn = false;
        }
      });

      $scope.logout = function() {
        userSvc.logout().then(function() {
          $state.go("root.home");
        });
      };
    }
  ]);
