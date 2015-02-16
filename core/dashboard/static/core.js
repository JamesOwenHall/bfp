// Entry point

window.onload = function() {

};

// Common utilities

// Performs a GET request.
function get(url, callback) {
  try {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
      if (xhr.readyState === 4) {
        callback(xhr.responseText);
      }
    };

    xhr.open("GET", url, true);
    xhr.send();
  } catch (e) {
    console.log(e);
  }
}
