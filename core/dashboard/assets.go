package dashboard

const (
	coreHtml = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Dashboard &middot; BFP</title>
    <link rel="stylesheet" href="/core.css">
  </head>
  <body>
    <div class="container">
      <header>
        <div id="logo">
          <h1>BFP</h1>
        </div>
      </header>

      <div class="divider">
        <h2>Dashboard</h2>
      </div>
      <section>
        <div class="grid">
          <div class="big-counter col-1-3">
            <p class="blocked-count">&nbsp;</p>
            <p>blocked values</p>
            <button id="refresh-button">Refresh</button>
          </div>
          <div class="col-2-3">
            <table>
              <tbody>
                <tr>
                  <td>Version</td>
                  <td>{{.Version}}</td>
                </tr>
                <tr>
                  <td>Listen type</td>
                  <td>{{.ListenType}}</td>
                </tr>
                <tr>
                  <td>Listen address</td>
                  <td>{{.ListenAddress}}</td>
                </tr>
                <tr>
                  <td>Number of directions</td>
                  <td>{{len .Directions}}</td>
                </tr>
                <tr>
                  <td>24-hour activity</td>
                  <td><span data-activity>0</span> hits</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <div class="divider">
        <h2>Directions</h2>
      </div>
      {{range .Directions}}
        <section data-direction="{{.Name}}">
          <h3>{{.Name}}</h3>
          <div class="grid">
            <div class="col-1-2">
              <table>
                <tbody>
                  <tr>
                    <td>Type</td>
                    <td>{{.Store.Type}}</td>
                  </tr>
                  <tr>
                    <td>Clean up time</td>
                    <td>{{.CleanUpTime}} sec</td>
                  </tr>
                  <tr>
                    <td>Threshold</td>
                    <td>{{.MaxHits}} hits</td>
                  </tr>
                  <tr>
                    <td>Window size</td>
                    <td>{{.WindowSize}} sec</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="col-1-2">
              <h4>Blocked Values</h4>
              <ul class="blocked-value-list"></ul>
            </div>
          </div>
        </section>
      {{end}}
    </div>
    <script src="/core.js"></script>
  </body>
</html>`

	coreCss = `
/* Page-wide */

*, *:before, *:after {
  -webkit-box-sizing: border-box;
  -moz-box-sizing: border-box;
  box-sizing: border-box;
}

:root {
  font-family: "Helvetica Neue", Arial, sans-serif;
  font-size: 14px;
}

body {
  background: hsl(0, 0%, 87%);
  margin: 0;
}

.container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  background: white;
}

h1, h2, h3, h4, h5, h6, p {
  margin: 0;
}

/* Grid */

.grid:after {
  content: "";
  display: table;
  clear: both;
}

[class*='col-'] {
  float: left;
  padding-right: 20px;
}
[class*='col-']:last-of-type {
  padding-right: 0;
}

.col-1-1 {
  width: 100%;
}
.col-1-2 {
  width: 50%;
}
.col-1-3 {
  width: 33%;
}
.col-2-3 {
  width: 67%;
}

@media (max-width: 768px) {
  [class*='col-'] {
    margin-top: 20px;
    width: 100%;
  }
  [class*='col-']:nth-child(1) {
    margin-top: 0;
  }
}

/* Tables */

table {
  border: 1px solid #CCC;
  border-collapse: collapse;
  width: 100%;
}

td {
  border: 1px solid #CCC;
  padding: 5px 8px;
}

td:nth-child(1) {
  font-weight: 500;
}

/* Header */

header {
  background: hsl(0, 0%, 33%);
  margin: 0;
}

header h1 {
  color: white;
  font-size: 1.7rem;
  padding: 30px 20px;
  text-shadow: 0 0 3px rgba(0, 0, 0, 0.5);
}

#logo {
  background: hsl(207, 50%, 48%);
  display: inline-block;
}

/* Sections */

section {
  padding: 20px;
}

.divider {
  background: hsl(0, 0%, 94%);
  margin: 0;
  padding: 12px 20px;
}

.divider h2 {
  font-size: 1rem;
}

h3 {
  color: hsl(0, 0%, 33%);
  text-transform: uppercase;
  padding-bottom: 0.5rem;
}

/* Big counter */

.big-counter {
  text-align: center;
}

.big-counter p {
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0;
}

.big-counter .blocked-count {
  font-size: 5rem;
}

/* Directions */

ul.blocked-value-list {
  max-height: 120px;
  overflow: scroll;
}
`

	coreJs = `
// The JSON response from the server.
var histData;

// Entry point
window.onload = function() {
  reloadData();
};

function reloadData() {
  get("/history", function(raw) {
    histData = JSON.parse(raw);

    updateBigCounter();
    update24hActivity();
    updateBlockedValueList();
  });
}

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

  histData.Directions.forEach(function(direction) {
    count += direction["blocked-values"].length
  });

  counter.innerText = count;
}

function update24hActivity() {
  var totalHits = histData.TotalHits;
  var td = document.querySelector("[data-activity]");
  td.innerText = totalHits;
}

document.getElementById("refresh-button").addEventListener("click", function() {
  reloadData();
});

// Blocked value list

function updateBlockedValueList() {
  histData.Directions.forEach(function(direction) {
    // Create the new list
    var fragment = document.createDocumentFragment();
    direction["blocked-values"].forEach(function(value) {
      var duration = histData.Clock - value.Since;
      var text = '“' + value.Value + '” for ' + duration + " seconds";

      var li = document.createElement("li");
      li.appendChild(document.createTextNode(text));
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
`
)
