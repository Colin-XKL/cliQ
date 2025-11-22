export function OpenUrlUseSystemBrowser(url: string) {
  window?.runtime?.BrowserOpenURL(url);
}