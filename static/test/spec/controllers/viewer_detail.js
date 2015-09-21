'use strict';

describe('Controller: ViewerDetailCtrl', function () {

  // load the controller's module
  beforeEach(module('scribbleApp'));

  var ViewerDetailCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    ViewerDetailCtrl = $controller('ViewerDetailCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(ViewerDetailCtrl.awesomeThings.length).toBe(3);
  });
});
