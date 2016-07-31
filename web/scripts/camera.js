(function() {

  var FRAME_RATE = 10;

  function Camera() {
    this._video = null;
    this.onError = null;
    this.onStart = null;
  }

  Camera.prototype.start = function() {
    getUserMedia(function(err, stream) {
      if (err !== null) {
        this.onError(err);
      } else {
        this._startWithStream(stream);
      }
    }.bind(this));
  };

  Camera.prototype.currentFrame = function() {
    var dims = this.outputDimensions();
    var canvas = document.createElement('canvas');
    canvas.width = dims.width;
    canvas.height = dims.height;
    var ctx = canvas.getContext('2d');
    ctx.drawImage(this._video, 0, 0, dims.width, dims.height);
    return canvas;
  };

  Camera.prototype.outputDimensions = function() {
    var width = this._video.videoWidth;
    var height = this._video.videoHeight;
    var scaleFactor = 192 / height;
    return {
      width: Math.round(width * scaleFactor),
      height: Math.round(height * scaleFactor)
    };
  };

  Camera.prototype._startWithStream = function(stream) {
    this._video = document.createElement('video');
    this._video.src = window.URL.createObjectURL(stream);
    this._video.play();

    // On chrome, onloadedmetadata will never be called, so we
    // use a timeout to start emitting frames anyway.
    var loadInterval;
    loadInterval = setInterval(function() {
      if (this._video.videoWidth) {
        clearInterval(loadInterval);
        this._video.onloadedmetadata = null;
        this.onStart();
      }
    }.bind(this), 500);

    this._video.onloadedmetadata = function() {
      clearInterval(loadInterval);
      this.onStart();
    }.bind(this);
  };

  window.app.Camera = Camera;

  function getUserMedia(cb) {
    var gum = (navigator.getUserMedia || navigator.webkitGetUserMedia ||
      navigator.mozGetUserMedia || navigator.msGetUserMedia);
    if (!gum) {
      setTimeout(function() {
        cb('getUserMedia() is not available.', null);
      }, 10);
      return;
    }
    gum.call(navigator, {audio: false, video: true},
      function(stream) {
        cb(null, stream);
      },
      function(err) {
        cb(err, null);
      }
    );
  }

})();
