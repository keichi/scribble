"use strict";

/**
 * @ngdoc function
 * @name scribbleApp.controller:LoginCtrl
 * @description
 * # LoginCtrl
 * Controller of the scribbleApp
 */
angular.module("scribbleApp")
  .controller("LoginCtrl", ["UserService", "$state", "ModalHelperService",
    function (userSvc, $state, modalSvc) {
      this.login = function() {
        userSvc
          .login(this.email, this.password)
          .then(function() {
            $state.go("root.home");
          }, function(resp) {
            if (resp.data) {
              modalSvc.alert("Login Error", resp.data.message);
            } else {
              modalSvc.alert("Login Error", "Unknown error.");
            }
          });
      };
    }
  ]);
