(function() {

  var currentWorker = null;

  window.app.detect = function(width, height, buffer, cb) {
    if (currentWorker === null) {
      currentWorker = new Worker('../detector/detector.js');
      currentWorker.postMessage(window.app.cascade);
    }
    currentWorker.onmessage = function(e) {
      cb(e.data);
    };
    currentWorker.postMessage([width, height, buffer]);
  };

})();
