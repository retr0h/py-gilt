***********
Development
***********

* Please read the :ref:`contributing` guidelines.

Branches
========

* The ``master`` branch is stable.  Major changes should be performed
  elsewhere.

Release Engineering
===================

Pre-release
-----------

* Edit the :ref:`changelog`.
* Follow the :ref:`testing` steps.

Release
-------

Gilt follows `Semantic Versioning`_.

.. _`Semantic Versioning`: http://semver.org

Tag the release and push to github.com
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: bash

    $ git tag 2.0.0
    $ git push --tags

Upload to `PyPI`_
^^^^^^^^^^^^^^^^^

* Upload to  `PyPI`_.

    .. code-block:: bash

        $ tox -e build-dists
        $ tox -e publish-dists

Post-release
------------

* Comment/close any relevant `Issues`_.

Roadmap
=======

* See `Issues`_ on Github.com.

.. _`PyPI`: https://pypi.python.org/pypi/python-gilt
.. _`ISSUES`: https://github.com/metacloud/gilt/issues
