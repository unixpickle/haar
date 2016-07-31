(function() {

  window.app = {};

  var camera = null;
  var canvas = null;

  function initialize() {
    canvas = document.getElementById('video-cell');
    camera = new window.app.Camera();
    camera.onError = function(e) {
      alert('Failed to configure camera: ' + e);
    };
    camera.onStart = function() {
      canvas.width = camera.outputDimensions().width;
      canvas.height = camera.outputDimensions().height;
      setInterval(function() {
        var ctx = canvas.getContext('2d');
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.drawImage(camera.currentFrame(), 0, 0);
      }, 100);
    };
    camera.start();
  }

  window.addEventListener('load', initialize);

})();
