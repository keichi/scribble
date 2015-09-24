/* global key */
"use strict";

/**
 * @ngdoc service
 * @name scribbleApp.ModalHelperService
 * @description
 * # ModalHelperService
 * Service in the scribbleApp.
 */
angular.module("scribbleApp")

  .service("ModalHelperService", function ($modal) {
    var service = this;

    service.alert = function(title, message) {
      var modalInstance = $modal.open({
        animation: true,
        templateUrl: "views/alert_modal.html",
        controller: "AlertModalInstanceCtrl",
        size: "sm",
        resolve: {
          title: function() {
            return title;
          },
          message: function() {
            return message;
          }
        }
      });

      return modalInstance.result;
    };
  })

  .controller("AlertModalInstanceCtrl", ["$scope", "$modalInstance", "title", "message",
    function($scope, $modalInstance, title, message) {
      $scope.title = title;
      $scope.message = message;

      var prevKeyScope = key.getScope();
      key.setScope("modal");

      $scope.ok = function () {
        $modalInstance.close(true);
      };

      $scope.cancel = function () {
        $modalInstance.dismiss("cancel");
      };

      $scope.$on("modal.closing", function() {
        key.setScope(prevKeyScope);
      });
    }
  ]);
