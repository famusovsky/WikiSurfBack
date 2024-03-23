document.currentScript = document.currentScript || (function() {
    var scripts = document.getElementsByTagName('script');
    return scripts[scripts.length - 1];
})();

let id = document.currentScript.getAttribute('rid');
let start = document.currentScript.getAttribute('start');
let finish = document.currentScript.getAttribute('finish');

chrome.runtime.sendMessage({
    action: 'wikiSurfBegin',
    wikiSurfStart: start,
    wikiSurfFinish: finish,
    wikiSurfId: id});

chrome.tabs.create({active: true, url: start});