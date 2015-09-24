'use strict';

describe('Service: ModalHelperService', function () {

  // load the service's module
  beforeEach(module('scribbleApp'));

  // instantiate service
  var ModalHelperService;
  beforeEach(inject(function (_ModalHelperService_) {
    ModalHelperService = _ModalHelperService_;
  }));

  it('should do something', function () {
    expect(!!ModalHelperService).toBe(true);
  });

});
