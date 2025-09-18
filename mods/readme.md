# hyprmods

This is a small tool to generate [Hyprland](https://hypr.land) binds which allow
modifier keys to be entered as sequences, alongside the (typical) chords. It
does so by generating a submap for each used modifier key combination, alongside
the release binds necessary to enter and exit them via the modifier keys.

## Usage

Simply add the following to the end of your `hyprland.conf`:

```
exec = hyprmods [MODS[=alias]...]
```

The tool will generate a submap for each listed modifier combination (using the
normal Hyprland syntax), using either the modifiers as listed or the provided
alias as the submap name.
