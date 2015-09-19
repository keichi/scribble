'use strict';

/**
 * @ngdoc service
 * @name scribbleApp.UserService
 * @description
 * # UserService
 * Service in the scribbleApp.
 */
angular.module('scribbleApp')
  .service('UserService', ['$http', 'ipCookie', 'API_ROOT',
    function ($http, ipCookie, API_ROOT) {
      var service = this;

      service.user = {};
      service.isLoggedIn = false;

      service.login = function(email, password) {
        return $http.post(API_ROOT + "/auth/login", {
          email: email,
          password: password
        }).then(function(resp) {
          ipCookie("token", resp.data.token, {
            expires: resp.data.sessionPeriod,
            expirationUnit: "milliseconds"
          });
          service.isLoggedIn = true;
        });
      };

      service.logout = function() {
        return $http.post(API_ROOT + "/auth/logout", {});
      };
  }]);
