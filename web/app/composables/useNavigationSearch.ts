export function useNavigationSearch() {
  const open = useState<boolean>("navigation-global-search-open", () => false);

  function openSearch() {
    open.value = true;
  }

  function closeSearch() {
    open.value = false;
  }

  return { open, openSearch, closeSearch };
}
