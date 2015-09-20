"use strict";

/**
 * @ngdoc service
 * @name scribbleApp.UserService
 * @description
 * # UserService
 * Service in the scribbleApp.
 */
angular.module("scribbleApp")

  .service("UserService", ["$http", "ipCookie", "API_ROOT",
    function ($http, ipCookie, API_ROOT) {
      var service = this;
      var isLoggedIn = false;

      service.login = function(email, password) {
        return $http.post(API_ROOT + "/auth/login", {
          email: email,
          password: password
        }).then(function(resp) {
          service.createSession(resp.data.token, resp.data.sessionPeriod);
          return resp;
        });
      };

      service.logout = function() {
        return $http.post(API_ROOT + "/auth/logout", {}).then(function () {
          service.invalidateSession();
        });
      };

      // Save specified session to our browser and change login state
      service.createSession = function(token, sessionPeriod) {
        ipCookie("token", token, {
          expires: sessionPeriod,
          expirationUnit: "milliseconds"
        });
        isLoggedIn = true;
      };

      // Delete session from our browser and change login state
      service.invalidateSession = function() {
        ipCookie.remove("token");
        isLoggedIn = false;
      };

      // Check current login state
      service.checkLoginStatus = function () {
        return $http.get(API_ROOT + "/auth", {}).then(function() {
          return true;
        }, function() {
          service.invalidateSession();
          return false;
        });
      };

      service.getSessionToken = function() {
        return ipCookie("token") || "";
      };
  }])

  .config(["$httpProvider", function($httpProvider) {
    $httpProvider.interceptors.push(["$q", "$injector",
      function($q, $injector) {
        return {
          request: function(config) {
            var userSvc = $injector.get("UserService");

            config.headers["X-Scribble-Session"] = userSvc.getSessionToken();

            return config;
          },
          responseError: function(rejection) {
            var userSvc = $injector.get("UserService");

            if (rejection.status === 500 && rejection.data) {
              if (rejection.data.message === "Session has expired" ||
                rejection.data.message === "not logged in") {
                userSvc.invalidateSession();
              }
            }
            return $q.reject(rejection);
          }
        };
      }
    ]);
  }]);
