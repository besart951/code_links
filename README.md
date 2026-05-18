# code_links

## Repo mit Submodules klonen

```bash id="7ax48j"
git clone --recurse-submodules https://github.com/besart951/code_links.git
cd code_links
```

## Falls ohne `--recurse-submodules` geklont wurde

```bash id="3efp8w"
git submodule update --init --recursive
```
