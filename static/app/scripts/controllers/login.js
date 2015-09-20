"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:LoginCtrl
 * @description
 * # LoginCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("LoginCtrl", ["UserService", "$state", function (userSvc, $state) {
    this.login = function() {
      userSvc
        .login(this.email, this.password)
        .then(function() {
          $state.go("home");
        }, function(resp) {
          if (resp.data) {
            window.alert(resp.data.message);
          } else {
            window.alert("unknown error");
          }
        });
    };
  }]);
