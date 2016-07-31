(function() {

  window.app = {};

  var camera = null;
  var canvas = null;
  var currentMatches = null;

  function initialize() {
    canvas = document.getElementById('video-cell');
    camera = new window.app.Camera();
    camera.onError = function(e) {
      alert('Failed to configure camera: ' + e);
    };
    camera.onStart = function() {
      canvas.width = camera.outputDimensions().width;
      canvas.height = camera.outputDimensions().height;
      var recognizing = false;
      setInterval(function() {
        var frame = camera.currentFrame();
        var ctx = canvas.getContext('2d');
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.drawImage(frame, 0, 0);

        if (currentMatches !== null) {
          ctx.strokeStyle = '#ff0000';
          for (var i = 0, len = currentMatches.length; i < len; ++i) {
            var match = currentMatches[i];
            ctx.strokeRect(match.X, match.Y, match.Width, match.Height);
          }
        }

        if (recognizing) {
          return;
        }
        recognizing = true;
        var matches = detectFacesInCanvas(frame, function(matches) {
          recognizing = false;
          currentMatches = matches;
        });
      }, 100);
    };
    camera.start();
  }

  function detectFacesInCanvas(canvas, cb) {
    var ctx = canvas.getContext('2d');
    var imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    var data = imageData.data;
    var floatData = [];
    var idx = 0;
    for (var i = 0, len = data.length/4; i < len; ++i) {
      floatData[i] = (data[idx]+data[idx+1]+data[idx+2]) / 765;
      idx += 4;
    }
    window.app.detect(canvas.width, canvas.height, floatData, cb);
  }

  window.addEventListener('load', initialize);

})();
