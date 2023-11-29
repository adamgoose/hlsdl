chrome.webRequest.onBeforeRequest.addListener(
  function(details) {
    console.log('URL:', details.url);
    return { cancel: false };
  },
  { urls: ["<all_urls>"] },
  ["blocking"]
);
