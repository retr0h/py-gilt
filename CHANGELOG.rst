*******
History
*******

1.2.3.
===

* Add bitbucket support.
* Add CI.
* Correct hard coded paths in tests.
* Remove black skip-string-normalization.
* Add coverage.
* Switched to setuptools vs pbr.
* Switched to black off yapf.
* Added hacking for extra lint.
* Corrected docstrings.

1.2.2
=====

* Cleanup dst prior to checkout-index.
* Fix ambigous variable usage.
* Fixed deprecation of getfuncargvalue.
* Remove use of deprecated pytest.config.
* Fixed test failure with GILT_CACHE_DIRECTORY.
* Fix tox failure due to use of relative path with --cov.

1.2.1
=====

* Use proper package name for pbr.

1.2
===

* Add option to override gilt's default cache dir.

1.1
===

* Add support for running commands after a sync has happened.
* Added python 3 support.
* Suppress doc building warnings.
* Only fetch when branch/tag/commit not on local repo.
* Determine branch in a smarter way.

1.0
===

* Initial release.
