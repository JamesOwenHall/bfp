// The JSON response from the server.
var histData;

// Entry point
window.onload = function() {
  get("/history", function(raw) {
    histData = JSON.parse(raw);

    updateBigCounter();
    update24hActivity();
    updateBlockedValueList();
  });
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

// Removes all child nodes from a DOM node.
function removeAllChildren(parent) {
  while (parent.firstChild) {
    parent.removeChild(parent.firstChild);
  }
}

// Big counter

function updateBigCounter() {
  var counter = document.querySelector(".big-counter .blocked-count");
  var count = 0;

  histData.forEach(function(direction) {
    count += direction["blocked-values"].length
  });

  counter.innerText = count;
}

function update24hActivity() {
  var totalHits = 0;

  histData.forEach(function(direction) {
    direction["long-history"].forEach(function(count) {
      totalHits += count;
    });
  });

  var td = document.querySelector("[data-activity]");
  td.innerText = totalHits;
}

// Blocked value list

function updateBlockedValueList() {
  histData.forEach(function(direction) {
    // Create the new list
    var fragment = document.createDocumentFragment();
    direction["blocked-values"].forEach(function(value) {
      var li = document.createElement("li");
      li.appendChild(document.createTextNode(value));
      fragment.appendChild(li);
    });

    // Remove the existing list
    var selector = "[data-direction='"+direction.name+"'] ul";
    var list = document.querySelector(selector);
    removeAllChildren(list);

    // Add the new one
    list.appendChild(fragment);
  });
}
