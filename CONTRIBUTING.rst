************
Contributing
************

* We are interested in various different kinds of improvement for Gilt;
  please feel free to raise an `Issue`_ if you would like to work on something
  major to ensure efficient collaboration and avoid duplicate effort.
* Create a topic branch from where you want to base your work.
* Check for unnecessary whitespace with ``git diff --check`` before committing.
* Make sure you have added tests for your changes.
* Run all the tests to ensure nothing else was accidentally broken.
* Reformat the code by following the formatting section below.
* Submit a pull request.

Testing
=======

Dependencies
------------

Install the test framework `Tox`_.

.. code-block:: bash

    $ pip install tox

Unit
----

Unit tests are invoked by `Tox`_.

.. code-block:: bash

    $ tox

Formatting
==========

The formatting is done using `Black`_.

From the root for the project, run:

.. code-block:: bash

    $ tox -e format

.. _`Black`: https://github.com/psf/black
.. _`Tox`: https://tox.readthedocs.io/en/latest
.. _`Issue`: https://github.com/metacloud/gilt/issues
