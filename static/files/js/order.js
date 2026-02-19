(function() {
  var page = document.getElementById("order-page");
  var catalog = JSON.parse(page.dataset.catalog);
  var members = JSON.parse(page.dataset.members);
  var club = page.dataset.club;
  var clubs = ["Unknown", "Parabool", "Gladiators", "Calamari"];
  var clubIndex = clubs.indexOf(club);
  if (clubIndex < 1) clubIndex = 2; // default Gladiators

  var categories = catalog.categories || [];
  var items = catalog.items || [];
  var cart = []; // [{item, amount}]
  var selectedCategoryID = 0;

  function priceForClub(item) {
    switch (clubIndex) {
      case 1: return item.price_parabool;
      case 2: return item.price_gladiators;
      case 3: return item.price_calamari;
      default: return 0;
    }
  }

  function formatPrice(cents) {
    var sign = cents < 0 ? "-" : "";
    cents = Math.abs(cents);
    return sign + "€" + (cents / 100).toFixed(2);
  }

  function visibleCategories() {
    var catHasItems = {};
    items.forEach(function(item) {
      if (priceForClub(item) !== 0) catHasItems[item.category_id] = true;
    });
    return categories.filter(function(c) { return catHasItems[c.id]; });
  }

  function visibleItems() {
    return items.filter(function(item) {
      return item.category_id === selectedCategoryID && priceForClub(item) !== 0;
    });
  }

  function calcTotal() {
    var total = 0;
    cart.forEach(function(line) {
      total += priceForClub(line.item) * line.amount;
    });
    return total;
  }

  function renderCategories() {
    var el = document.getElementById("category-list");
    el.innerHTML = "";
    visibleCategories().forEach(function(cat) {
      var btn = document.createElement("button");
      btn.className = "responsive extra small-margin" + (cat.id === selectedCategoryID ? " secondary" : "");
      btn.innerHTML = '<div style="text-overflow:ellipsis;white-space:nowrap;overflow:hidden;">' + esc(cat.name) + '</div>';
      btn.onclick = function() { selectedCategoryID = cat.id; render(); };
      el.appendChild(btn);
    });
  }

  function renderItems() {
    var el = document.getElementById("item-list");
    el.innerHTML = "";
    visibleItems().forEach(function(item) {
      var btn = document.createElement("button");
      btn.className = "responsive extra small-margin";
      btn.innerHTML = '<div class="row no-wrap">' +
        '<div class="col">' + esc(item.name) + '</div>' +
        '<div class="col min">' + formatPrice(priceForClub(item)) + '</div>' +
        '</div>';
      btn.onclick = function() { addToCart(item); };
      el.appendChild(btn);
    });
  }

  function renderOverview() {
    var el = document.getElementById("overview");
    el.innerHTML = "";
    cart.forEach(function(line) {
      var price = priceForClub(line.item);
      var cls = price === 0 ? "error" : "";
      var div = document.createElement("div");
      div.className = "small-padding " + cls;
      div.innerHTML =
        '<div class="row no-wrap large-text no-club">' +
          '<div class="col min middle-align">' +
            '<button class="circle flat left-round no-margin" data-action="remove" data-id="' + line.item.id + '"><i>remove</i></button>' +
          '</div>' +
          '<div class="col min middle-align" style="text-align:center;"><span class="bold" style="width:20px">' + line.amount + '</span></div>' +
          '<div class="col min middle-align">' +
            '<button class="circle flat right-round no-margin" data-action="add" data-id="' + line.item.id + '"><i>add</i></button>' +
          '</div>' +
          '<div class="col max middle-align"><span>' + esc(line.item.name) + '</span></div>' +
          '<div class="col min middle-align">' + formatPrice(price * line.amount) + '</div>' +
          '<div class="col min middle-align">' +
            '<button class="circle flat error" data-action="delete" data-id="' + line.item.id + '"><i>delete</i></button>' +
          '</div>' +
        '</div>';
      el.appendChild(div);
    });

    el.querySelectorAll("[data-action]").forEach(function(btn) {
      btn.onclick = function() {
        var id = parseInt(btn.dataset.id);
        var item = findItemById(id);
        if (!item) return;
        if (btn.dataset.action === "add") addToCart(item);
        else if (btn.dataset.action === "remove") removeFromCart(item);
        else if (btn.dataset.action === "delete") deleteFromCart(item);
      };
    });
  }

  function renderSummary() {
    var el = document.getElementById("summary");
    var disabled = cart.length === 0 ? "disabled" : "";
    el.innerHTML =
      '<div style="display:grid;grid-template-columns:1fr 1fr;grid-template-rows:1fr 1fr;grid-template-areas:\'total total\' \'btn1 btn2\';">' +
        '<h3 style="grid-area:total;justify-self:right;">Total ' + formatPrice(calcTotal()) + '</h3>' +
        '<div style="justify-self:center;">' +
          '<button id="btn-anon" class="extra" style="width:100px;" ' + disabled + '><i>payment</i><label class="primary">Anonymous</label></button>' +
        '</div>' +
        '<div style="justify-self:center;">' +
          '<button id="btn-member" class="extra" style="width:100px;" ' + disabled + '><i>person</i><label class="primary">Member</label></button>' +
        '</div>' +
      '</div>';
    var btnAnon = document.getElementById("btn-anon");
    var btnMember = document.getElementById("btn-member");
    if (btnAnon) btnAnon.onclick = showAnonPayment;
    if (btnMember) btnMember.onclick = showMemberPicker;
  }

  var logoMap = {
    "Parabool": "/static/logos/parabool.jpg",
    "Gladiators": "/static/logos/gladiators.jpg",
    "Calamari": "/static/logos/calamari.jpg"
  };

  function renderToggler() {
    var el = document.getElementById("club-toggler");
    var name = clubs[clubIndex];
    var logo = logoMap[name] || "";
    var bgStyle = logo
      ? "background-color:white;background-position:center;background-repeat:no-repeat;" +
        "background-image:url(" + logo + ");background-size:120px;" +
        "height:120px;width:120px;padding:24px;border-radius:84px;border:none;"
      : "width:120px;height:120px;";
    el.innerHTML = '<div class="row no-wrap"><div class="col max"></div>' +
      '<div class="col min"><div id="club-toggle-btn" style="cursor:pointer;' + bgStyle + '"></div></div>' +
      '<div class="col max"></div></div>';
    document.getElementById("club-toggle-btn").onclick = toggleClub;
  }

  function render() {
    page.className = "full-height page-content " + clubs[clubIndex];
    renderCategories();
    renderItems();
    renderOverview();
    renderSummary();
    renderToggler();
  }

  function findItemById(id) {
    for (var i = 0; i < items.length; i++) {
      if (items[i].id === id) return items[i];
    }
    return null;
  }

  function addToCart(item) {
    for (var i = 0; i < cart.length; i++) {
      if (cart[i].item.id === item.id) { cart[i].amount++; render(); return; }
    }
    cart.push({item: item, amount: 1});
    render();
  }

  function removeFromCart(item) {
    for (var i = 0; i < cart.length; i++) {
      if (cart[i].item.id === item.id) {
        if (cart[i].amount <= 1) { cart.splice(i, 1); } else { cart[i].amount--; }
        render(); return;
      }
    }
  }

  function deleteFromCart(item) {
    cart = cart.filter(function(line) { return line.item.id !== item.id; });
    render();
  }

  function toggleClub() {
    if (clubIndex === 2) clubIndex = 1;      // Gladiators -> Parabool
    else if (clubIndex === 1) clubIndex = 3;  // Parabool -> Calamari
    else clubIndex = 2;                       // Calamari -> Gladiators
    render();
  }

  function esc(s) {
    var d = document.createElement("div");
    d.textContent = s;
    return d.innerHTML;
  }

  // --- Member picker ---

  function showMemberPicker() {
    document.getElementById("member-modal").style.display = "block";
    document.getElementById("member-search").value = "";
    renderMemberList("");
    setTimeout(function() { document.getElementById("member-search").focus(); }, 50);
  }

  window.closeMemberModal = function() {
    document.getElementById("member-modal").style.display = "none";
  };

  window.filterMembers = function() {
    var q = document.getElementById("member-search").value.toLowerCase();
    renderMemberList(q);
  };

  function renderMemberList(query) {
    var el = document.getElementById("member-list");
    el.innerHTML = "";
    var filtered = members.filter(function(m) {
      return m.club === clubs[clubIndex] && (query === "" || m.name.toLowerCase().indexOf(query) >= 0);
    });
    filtered.sort(function(a, b) {
      return new Date(b.last_order) - new Date(a.last_order);
    });
    filtered.forEach(function(m) {
      var btn = document.createElement("button");
      btn.className = "responsive extra small-margin " + m.club;
      btn.innerHTML = '<div style="text-overflow:ellipsis;white-space:nowrap;overflow:hidden;">' + esc(m.name) + '</div>';
      btn.onclick = function() { selectMember(m); };
      el.appendChild(btn);
    });
  }

  // --- Payment modals ---

  function showPaymentModal(html) {
    document.getElementById("payment-content").innerHTML = html;
    document.getElementById("payment-modal").style.display = "block";
  }

  function closePaymentModal() {
    document.getElementById("payment-modal").style.display = "none";
  }

  document.getElementById("payment-back-btn").onclick = closePaymentModal;
  document.getElementById("payment-backdrop").onclick = closePaymentModal;

  function showAnonPayment() {
    showPaymentModal(
      '<h4>Anonymous</h4>' +
      '<div style="font-size:2em;margin-bottom:15px;" class="row no-wrap">' +
        '<span class="col">Price</span>' +
        '<span class="col min">' + formatPrice(calcTotal()) + '</span>' +
      '</div>' +
      '<p>Pay by PIN</p>' +
      '<button id="pay-anon-btn" class="responsive large">Paid</button>' +
      '<div id="pay-error" style="display:none;" class="error" style="margin-top:10px;">Unable to place order.</div>'
    );
    document.getElementById("pay-anon-btn").onclick = function() { placeOrder(0); };
  }

  function selectMember(member) {
    closeMemberModal();
    showPaymentModal(
      '<h4>' + esc(member.name) + '</h4>' +
      '<div id="member-debt-loading">Loading...</div>' +
      '<div id="member-debt-info" style="display:none;">' +
        '<div class="row no-wrap large-text"><span class="col">Current Bill</span><span class="col min" id="member-debt"></span></div>' +
        '<div style="font-size:2em;margin-bottom:15px;" class="row no-wrap"><span class="col">Price</span><span class="col min">' + formatPrice(calcTotal()) + '</span></div>' +
        '<button id="pay-member-btn" class="responsive large">Add to Bill</button>' +
      '</div>' +
      '<div id="pay-error" style="display:none;margin-top:10px;" class="error">Unable to place order.</div>'
    );

    fetch("/api/member/" + member.id, {credentials: "same-origin"})
      .then(function(r) { return r.json(); })
      .then(function(data) {
        document.getElementById("member-debt-loading").style.display = "none";
        document.getElementById("member-debt").textContent = formatPrice(data.debt);
        document.getElementById("member-debt-info").style.display = "block";
        document.getElementById("pay-member-btn").onclick = function() { placeOrder(member.id); };
      })
      .catch(function() {
        document.getElementById("member-debt-loading").textContent = "Unable to load member details.";
      });
  }

  function placeOrder(memberID) {
    var orderLines = cart.map(function(line) {
      return {product: line.item, amount: line.amount};
    });
    var orderData = {
      club: clubs[clubIndex],
      memberId: memberID,
      contents: JSON.stringify(orderLines),
      price: calcTotal(),
      status: "Open"
    };

    var btn = document.getElementById(memberID === 0 ? "pay-anon-btn" : "pay-member-btn");
    if (btn) btn.disabled = true;

    fetch("/api/order", {
      method: "POST",
      credentials: "same-origin",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(orderData)
    }).then(function(r) {
      if (!r.ok) throw new Error("order failed");
      cart = [];
      closePaymentModal();
      render();
    }).catch(function() {
      var errEl = document.getElementById("pay-error");
      if (errEl) { errEl.style.display = "block"; }
      if (btn) btn.disabled = false;
    });
  }

  // --- Grid CSS ---

  var style = document.createElement("style");
  style.textContent =
    '#order-grid { height:100%;display:grid;grid-gap:5px;' +
    'grid-template-columns:30% 30% 40%;grid-template-rows:50px 1fr 200px;' +
    "grid-template-areas:'cat-title item-title over-title' 'categories items overview' 'toggler items payment';}" +
    '@media (max-width:900px) { #order-grid { grid-template-columns:1fr;' +
    'grid-template-rows:auto auto auto auto auto auto auto auto auto;' +
    "grid-template-areas:'cat-title' 'categories' 'item-title' 'items' 'over-title' 'overview' 'toggler' 'payment';}}";
  document.head.appendChild(style);

  render();
})();
