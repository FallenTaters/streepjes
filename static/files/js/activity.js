(function() {
  var INACTIVITY_TIMEOUT_MS = 15 * 60 * 1000;
  var WARNING_BEFORE_MS = 10 * 60 * 1000;
  var KEEPALIVE_INTERVAL_MS = 2 * 60 * 1000;

  var lastActivity = Date.now();
  var warningShown = false;
  var overlay = null;

  function resetActivity() {
    lastActivity = Date.now();
    if (warningShown) {
      hideWarning();
    }
  }

  function showWarning() {
    if (warningShown) return;
    warningShown = true;
    overlay = document.createElement("div");
    overlay.style.cssText =
      "position:fixed;top:0;left:0;width:100vw;height:100vh;" +
      "background:black;color:white;z-index:10000;" +
      "display:flex;align-items:center;justify-content:center;" +
      "font-size:2em;cursor:pointer;";
    overlay.textContent = "Session expiring soon. Click anywhere to stay logged in.";
    overlay.addEventListener("click", function() {
      resetActivity();
      sendKeepalive();
    });
    document.body.appendChild(overlay);
  }

  function hideWarning() {
    warningShown = false;
    if (overlay && overlay.parentNode) {
      overlay.parentNode.removeChild(overlay);
      overlay = null;
    }
  }

  function sendKeepalive() {
    fetch("/active", { method: "POST", credentials: "same-origin" }).catch(function() {});
  }

  function check() {
    var elapsed = Date.now() - lastActivity;
    if (elapsed >= INACTIVITY_TIMEOUT_MS) {
      window.location.href = "/logout";
      return;
    }
    if (elapsed >= INACTIVITY_TIMEOUT_MS - WARNING_BEFORE_MS) {
      showWarning();
    }
  }

  document.addEventListener("click", resetActivity);
  document.addEventListener("keypress", resetActivity);

  setInterval(check, 10000);
  setInterval(function() {
    if (Date.now() - lastActivity < INACTIVITY_TIMEOUT_MS - WARNING_BEFORE_MS) {
      sendKeepalive();
    }
  }, KEEPALIVE_INTERVAL_MS);
})();
