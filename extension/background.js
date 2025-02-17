async function getAllowedSites() {
  let { allowedSites } = await chrome.storage.local.get(['allowedSites'])

  return new Set(allowedSites || [])
}

async function setAllowedSites(allowedSites) {
  await chrome.storage.local.set({ allowedSites: Array.from(allowedSites) })
}

async function updateContextMenu(url) {
  if (!url) return

  var url = new URL(url);
  var domain = url.hostname;
  let allowedSites = await getAllowedSites()
  if (allowedSites.has(domain)) {
    chrome.contextMenus.update("toggle_m3u8_scraping", { enabled: true, title: "Disable Domain: " + domain });
  }
  else {
    chrome.contextMenus.update("toggle_m3u8_scraping", { enabled: true, title: "Enable Domain: " + domain });
  }
}

chrome.webRequest.onBeforeRequest.addListener(
  async function(details) {
    let allowedSites = await getAllowedSites()

    let url = new URL(details.initiator)
    if (!allowedSites.has(url.hostname)) {
      return { cancel: false };
    }

    url = new URL(details.url)
    if (!url.toString().includes('.m3u8')) {
      return { cancel: false };
    }

    let tab = await chrome.tabs.get(details.tabId)

    await chrome.scripting.executeScript({
      target: { tabId: details.tabId },
      func: (url, downloadUrl) => {
        console.log('HLSDL Detection!')
        console.log('URL: ', url)
        console.log('Download: ', downloadUrl)
        console.log('---------------------------')
      },
      args: [
        url.toString(),
        'http://totoro:8881/dl?url=' + encodeURIComponent(url.toString()) + '&filename=' + encodeURIComponent(tab.title)
      ]
    })

    console.log(details)

    return { cancel: false };
  },
  { urls: ["<all_urls>"] },
  []
);

// A generic onclick callback function.
chrome.contextMenus.onClicked.addListener(async event => {
  var url = new URL(event.pageUrl)
  var domain = url.hostname;

  switch (event.menuItemId) {
    case 'toggle_m3u8_scraping':
      let allowedSites = await getAllowedSites()
      if (allowedSites.has(domain)) {
        allowedSites.delete(domain)
      }
      else {
        allowedSites.add(domain)
      }

      await setAllowedSites(allowedSites)
      await updateContextMenu(event.pageUrl)
      break
  }
});

chrome.runtime.onInstalled.addListener(function() {
  chrome.contextMenus.create({
    title: 'Toggle m3u8 Scraping',
    contexts: ['page'],
    id: 'toggle_m3u8_scraping',
    enabled: false
  })
});


chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === 'complete') {
    updateContextMenu(tab.url);
  }
});

chrome.tabs.onActivated.addListener((activeInfo) => {
  chrome.tabs.get(activeInfo.tabId, (tab) => {
    updateContextMenu(tab.url);
  })
});
