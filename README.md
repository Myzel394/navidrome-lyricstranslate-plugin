# Lyricstranslate Lyrics Plugin for Navidrome

Scrapes lyrics from Lyricstranslate and provides them to your Navidrome instance. Scraping is not implemented yet; the plugin currently returns constant stub lyrics.

## Installation

1. Download `navidrome-lyricstranslate-plugin.ndp` from the [latest release](https://github.com/Myzel394/navidrome-lyricstranslate-plugin/releases/latest).
2. Copy it to your Navidrome plugins folder (default: `<navidrome-data-directory>/plugins/`).
3. Add `navidrome-lyricstranslate-plugin` to the lyrics priority list (e.g. using envs: `ND_LYRICSPRIORITY=other-lyric-provider,navidrome-lyricstranslate-plugin`)
4. In Navidrome, go to **Settings > Plugins > Navidrome Plugin** and toggle it on.

It's recommended to set this plugin's priority to the lowest position, as scraping is less reliable than using an API.

**Lyricstranslate scraping is not implemented yet. The plugin currently returns stub lyrics.**

## Reporting Issues

Before opening an [issue](https://github.com/Myzel394/navidrome-lyricstranslate-plugin/issues), grep your Navidrome logs and attach the matching lines:

```sh
grep navidrome-lyricstranslate-plugin /path/to/navidrome.log
```
