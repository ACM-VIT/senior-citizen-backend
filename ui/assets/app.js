(() => {
  const cfg = window.AUTH_UI_CONFIG || {};
  const kratosPublicUrl = (cfg.kratosPublicUrl || `${window.location.protocol}//${window.location.hostname}:4433`).replace(/\/+$/, "");
  const publicUrl = (cfg.publicUrl || window.location.origin).replace(/\/+$/, "");
  let sessionPromise = null;

  const path = window.location.pathname.replace(/\/+$/, "");
  if (path === "/ui") {
    window.location.replace("/ui/login");
    return;
  }

  const flowType = path.split("/").pop() || "login";
  const params = new URLSearchParams(window.location.search);

  const views = {
    login: {
      title: "Sign in to ACM-VIT",
      subtitle: "Access every ACM-VIT app with one account.",
      links: [
        { text: "Create an account", href: "/ui/registration" },
        { text: "Forgot password", href: "/ui/recovery" }
      ]
    },
    registration: {
      title: "Create your ACM-VIT account",
      subtitle: "Get started with a single sign-on profile.",
      links: [
        { text: "Already have an account", href: "/ui/login" }
      ]
    },
    settings: {
      title: "Account settings",
      subtitle: "Manage your ACM-VIT profile and security.",
      links: [
        { text: "Back to app", href: publicUrl }
      ]
    },
    recovery: {
      title: "Recover access",
      subtitle: "Send a recovery link to your ACM-VIT email.",
      links: [
        { text: "Back to sign in", href: "/ui/login" }
      ]
    },
    verification: {
      title: "Verify your email",
      subtitle: "Finish verifying your ACM-VIT email.",
      links: [
        { text: "Back to sign in", href: "/ui/login" }
      ]
    },
    error: {
      title: "Something went wrong",
      subtitle: "Review the details below and try again.",
      links: [
        { text: "Back to sign in", href: "/ui/login" }
      ]
    },
    logout: {
      title: "Signing out",
      subtitle: "Clearing your session.",
      links: []
    }
  };

  const view = views[flowType];
  if (!view) {
    setTitle("Unknown page");
    setSubtitle("This page does not exist.");
    showMessage("Invalid path: " + flowType, "error");
    return;
  }

  setTitle(view.title);
  setSubtitle(view.subtitle);
  renderLinks(view.links, params);
  setEnvHost(kratosPublicUrl);

  const flowId = params.get("flow");
  init().catch((err) => {
    showMessage("Something went wrong. " + err.message, "error");
  });

  async function init() {
    if (flowType === "logout") {
      renderStatus("Signing you out...");
      window.location.href = `${kratosPublicUrl}/self-service/logout/browser`;
      return;
    }

    if (flowType === "error") {
      await handleError(params);
      return;
    }

    const gated = await maybeShowSessionGate(flowType, params, kratosPublicUrl);
    if (gated) {
      return;
    }

    if (!flowId) {
      startFlow(flowType, params, kratosPublicUrl);
      return;
    }

    try {
      const flow = await fetchFlow(flowType, flowId, kratosPublicUrl);
      renderFlow(flow, flowType);
    } catch (err) {
      if (shouldRestartFlow(err, flowType, flowId)) {
        showMessage("Flow expired. Starting a fresh one...", "error");
        restartFlow(flowType, params, kratosPublicUrl, flowId);
        return;
      }
      showMessage("Could not load the flow. " + err.message, "error");
    }
  }

  function setTitle(text) {
    const el = document.getElementById("title");
    if (el) {
      el.textContent = text;
    }
  }

  function setSubtitle(text) {
    const el = document.getElementById("subtitle");
    if (el) {
      el.textContent = text;
    }
  }

  function setEnvHost(url) {
    const el = document.getElementById("env-host");
    if (!el) {
      return;
    }
    el.textContent = url.replace(/^https?:\/\//, "");
  }

  function renderStatus(text) {
    const container = document.getElementById("form-container");
    if (container) {
      container.innerHTML = `<div class="status">${escapeHTML(text)}</div>`;
    }
  }

  function renderLinks(links, params) {
    const container = document.getElementById("aux-links");
    if (!container) {
      return;
    }
    container.innerHTML = "";
    if (!links || links.length === 0) {
      return;
    }

    const list = document.createElement("div");
    list.className = "link-list";
    links.forEach((link) => {
      const a = document.createElement("a");
      a.href = withParams(link.href, params);
      a.textContent = link.text;
      list.appendChild(a);
    });
    container.appendChild(list);
  }

  function withParams(path, params) {
    const url = new URL(path, window.location.origin);
    ["return_to", "refresh"].forEach((key) => {
      if (params.has(key)) {
        url.searchParams.set(key, params.get(key));
      }
    });
    return url.pathname + url.search;
  }

  function startFlow(flowType, params, kratosPublicUrl) {
    const url = new URL(`${kratosPublicUrl}/self-service/${flowType}/browser`);
    params.forEach((value, key) => {
      if (key !== "flow") {
        url.searchParams.append(key, value);
      }
    });
    window.location.href = url.toString();
  }

  async function fetchFlow(flowType, flowId, kratosPublicUrl) {
    const url = `${kratosPublicUrl}/self-service/${flowType}/flows?id=${encodeURIComponent(flowId)}`;
    const resp = await fetch(url, {
      credentials: "include",
      headers: { "Accept": "application/json" }
    });
    if (!resp.ok) {
      let payload;
      try {
        payload = await resp.json();
      } catch (error) {
        payload = null;
      }
      const message = payload && payload.error && payload.error.message
        ? payload.error.message
        : `Flow request failed (${resp.status}).`;
      const err = new Error(message);
      err.status = resp.status;
      err.payload = payload;
      throw err;
    }
    return resp.json();
  }

  function shouldRestartFlow(err, flowType, flowId) {
    if (!err || !flowId) {
      return false;
    }
    const status = err.status;
    if (status !== 404 && status !== 410) {
      return false;
    }
    const key = `flow-restart:${flowType}:${flowId}`;
    if (window.sessionStorage.getItem(key)) {
      return false;
    }
    window.sessionStorage.setItem(key, "1");
    return true;
  }

  function restartFlow(flowType, params, kratosPublicUrl, flowId) {
    if (flowId) {
      params.delete("flow");
    }
    startFlow(flowType, params, kratosPublicUrl);
  }

  function shouldGateSession(flowType, params) {
    if (flowType !== "login" && flowType !== "registration") {
      return false;
    }
    return !isReauthRequested(params);
  }

  function isReauthRequested(params) {
    const refresh = params.get("refresh");
    const aal = params.get("aal");
    return refresh === "true" || refresh === "1" || Boolean(aal);
  }

  async function maybeShowSessionGate(flowType, params, kratosPublicUrl) {
    if (!shouldGateSession(flowType, params)) {
      return false;
    }
    const session = await getSession(kratosPublicUrl);
    if (!session) {
      return false;
    }
    renderSessionGate(session, params);
    return true;
  }

  function getSession(kratosPublicUrl) {
    if (!sessionPromise) {
      sessionPromise = (async () => {
        try {
          const resp = await fetch(`${kratosPublicUrl}/sessions/whoami`, {
            credentials: "include",
            headers: { "Accept": "application/json" }
          });
          if (!resp.ok) {
            return null;
          }
          const data = await resp.json();
          if (data && data.active === false) {
            return null;
          }
          return data;
        } catch (error) {
          return null;
        }
      })();
    }
    return sessionPromise;
  }

  function renderSessionGate(session, params) {
    setTitle("You're already signed in");
    setSubtitle("Continue to ACM-VIT apps or sign out to switch accounts.");
    renderLinks([], params);
    renderMessages([], "messages");

    const container = document.getElementById("form-container");
    if (!container) {
      return;
    }
    container.innerHTML = "";

    const status = document.createElement("div");
    status.className = "status";
    const identifier = getSessionIdentifier(session);
    status.textContent = identifier
      ? `Signed in as ${identifier}.`
      : "Your session is active.";
    container.appendChild(status);

    const actions = document.createElement("div");
    actions.className = "session-actions";

    const continueLink = document.createElement("a");
    continueLink.className = "btn primary";
    continueLink.href = resolveReturnTo(params);
    continueLink.textContent = "Continue";
    actions.appendChild(continueLink);

    const logoutLink = document.createElement("a");
    logoutLink.className = "btn";
    logoutLink.href = "/ui/logout";
    logoutLink.textContent = "Sign out";
    actions.appendChild(logoutLink);

    container.appendChild(actions);
  }

  function getSessionIdentifier(session) {
    const traits = session && session.identity && session.identity.traits;
    if (!traits) {
      return "";
    }
    return traits.email || traits.username || traits.name || traits.phone || "";
  }

  function resolveReturnTo(params) {
    return params.get("return_to") || publicUrl;
  }

  async function handleError(params) {
    const errId = params.get("id");
    const hydraError = params.get("error");
    const hydraDesc = params.get("error_description");
    if (!errId && (hydraError || hydraDesc)) {
      const message = [hydraError, hydraDesc].filter(Boolean).join(": ");
      showMessage(message, "error");
      return;
    }
    if (!errId) {
      showMessage("No error ID provided.", "error");
      return;
    }
    const url = `${kratosPublicUrl}/self-service/errors?id=${encodeURIComponent(errId)}`;
    const resp = await fetch(url, {
      credentials: "include",
      headers: { "Accept": "application/json" }
    });
    if (!resp.ok) {
      showMessage("Failed to load error details.", "error");
      return;
    }
    const data = await resp.json();
    if (data && data.error && data.error.message) {
      showMessage(data.error.message, "error");
    } else {
      showMessage("An unknown error occurred.", "error");
    }
  }

  function renderFlow(flow, flowType) {
    if (!flow || !flow.ui) {
      showMessage("Missing form data.", "error");
      return;
    }

    const container = document.getElementById("form-container");
    if (!container) {
      return;
    }
    container.innerHTML = "";

    renderMessages(flow.ui.messages || [], "messages");

    const form = document.createElement("form");
    form.className = "auth-form";
    form.method = flow.ui.method || "POST";
    form.action = flow.ui.action || "";
    form.autocomplete = "on";

    const nodes = flow.ui.nodes || [];
    if (nodes.length === 0) {
      const empty = document.createElement("div");
      empty.className = "status";
      empty.textContent = "No fields to render. Try refreshing the page.";
      container.appendChild(empty);
      return;
    }

    nodes.forEach((node) => {
      renderNode(node, form, flowType);
    });

    container.appendChild(form);
  }

  function renderNode(node, form, flowType) {
    if (!node) {
      return;
    }

    if (node.type === "text") {
      const text = (node.text && node.text.text) || (node.meta && node.meta.label && node.meta.label.text) || "";
      if (text) {
        const p = document.createElement("p");
        p.className = "form-text";
        p.textContent = text;
        form.appendChild(p);
      }
      return;
    }

    if (node.type !== "input" || !node.attributes) {
      return;
    }

    const attrs = node.attributes;
    const inputType = attrs.type || "text";

    if (inputType === "hidden") {
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = attrs.name || "";
      if (attrs.value !== undefined) {
        input.value = attrs.value;
      }
      form.appendChild(input);
      return;
    }

    if (inputType === "submit") {
      const button = document.createElement("button");
      button.type = "submit";
      button.className = "btn primary";
      if (attrs.name) {
        button.name = attrs.name;
      }
      if (attrs.value !== undefined) {
        button.value = attrs.value;
      }
      if (attrs.disabled) {
        button.disabled = true;
      }
      button.textContent = labelForNode(node) || attrs.value || "Continue";
      form.appendChild(button);
      return;
    }

    const field = document.createElement("div");
    field.className = "field";

    const labelText = labelForNode(node);
    if (labelText) {
      const label = document.createElement("label");
      const inputId = attrs.id || attrs.name || "field-" + Math.random().toString(36).slice(2);
      label.htmlFor = inputId;
      label.textContent = labelText;
      field.appendChild(label);

      const input = buildInput(attrs, inputType, inputId);
      field.appendChild(input);
      renderNodeMessages(node, field);
      form.appendChild(field);
      return;
    }

    const input = buildInput(attrs, inputType, attrs.id || "");
    field.appendChild(input);
    renderNodeMessages(node, field);
    form.appendChild(field);
  }

  function buildInput(attrs, inputType, id) {
    const input = document.createElement("input");
    input.type = inputType;
    if (attrs.name) {
      input.name = attrs.name;
    }
    if (id) {
      input.id = id;
    }
    if (attrs.value !== undefined && inputType !== "password") {
      input.value = attrs.value;
    }
    if (attrs.required) {
      input.required = true;
    }
    if (attrs.disabled) {
      input.disabled = true;
    }
    if (attrs.autocomplete) {
      input.autocomplete = attrs.autocomplete;
    }
    if (attrs.pattern) {
      input.pattern = attrs.pattern;
    }
    if (attrs.min_length) {
      input.minLength = attrs.min_length;
    }
    if (attrs.max_length) {
      input.maxLength = attrs.max_length;
    }
    if (!input.placeholder && attrs.name) {
      input.placeholder = prettify(attrs.name);
    }
    return input;
  }

  function labelForNode(node) {
    if (node.meta && node.meta.label && node.meta.label.text) {
      return node.meta.label.text;
    }
    return "";
  }

  function renderNodeMessages(node, container) {
    if (!node.messages || node.messages.length === 0) {
      return;
    }
    const wrapper = document.createElement("div");
    wrapper.className = "field-messages";
    node.messages.forEach((msg) => {
      const div = document.createElement("div");
      div.className = "msg msg--error";
      div.textContent = msg.text || "Invalid input.";
      wrapper.appendChild(div);
    });
    container.appendChild(wrapper);
  }

  function renderMessages(messages, containerId) {
    const container = document.getElementById(containerId);
    if (!container) {
      return;
    }
    container.innerHTML = "";
    if (!messages || messages.length === 0) {
      return;
    }
    messages.forEach((msg) => {
      const div = document.createElement("div");
      div.className = msg.type === "error" ? "msg msg--error" : "msg";
      div.textContent = msg.text || "";
      container.appendChild(div);
    });
  }

  function showMessage(text, type) {
    renderMessages([{ text, type }], "messages");
  }

  function escapeHTML(value) {
    return String(value)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/\"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  function prettify(value) {
    return value
      .replace(/_/g, " ")
      .replace(/\./g, " ")
      .replace(/\b\w/g, (c) => c.toUpperCase());
  }
})();
