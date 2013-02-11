# l10n_check

l10n_check is a small utility to validate one or multiple java property files _(.properties)_.  
It is (as the name suggests) intended to validate localization files.  
Thus, when multiple files are validated, missing keys and such are reported.

## Usage:
l10n_check [-v] [-nowarn] <file-name> [<file-name> ...]
the optional -v (verbose) switch causes l10n_check to print all key-value pairs.
the optional -nowarn switch causes  l10n_check to only print errors.

### Exit codes:
0 : Everything is ok (as far as l10n_check is concerned)
1 : Bad command-line parameters (no files supplied)
2 : There are validation errors
