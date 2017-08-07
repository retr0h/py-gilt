.. include:: ../../README.rst

Usage:
======

Overlay a remote repository into the destination provided.

.. code-block:: yaml
  :caption: gilt.yml

    - git: https://github.com/retr0h/ansible-etcd.git
      version: master
      dst: roles/retr0h.ansible-etcd/

.. code-block:: bash

    $ gilt overlay


Overlay files from a remote repository into the destinations provided.

.. code-block:: yaml
  :caption: gilt.yml

    - git: https://github.com/lorin/openstack-ansible-modules.git
      version: master
      files:
        - src: "*_manage"
          dst: library/
        - src: nova_quota
          dst: library/
        - src: neutron_router
          dst: library/neutron_router.py

.. code-block:: bash

    $ gilt overlay

Overlay a directory from a remote repository into the destination provided.

.. code-block:: yaml
  :caption: gilt.yml

    - git: https://github.com/blueboxgroup/ursula.git
      version: master
      files:
        - src: roles/logging
          dst: roles/blueboxgroup.logging/

.. code-block:: bash

    $ gilt overlay

Optionally, override gilt's cache location (defaults to ~/.gilt):

.. code-block:: bash

    $ export GILT_CACHE_DIRECTORY=~/my-gilt-cache

Overlay files and a directory and run post-overlay commands.

.. code-block:: yaml
  :caption: gilt.yml

    - git: https://github.com/example/subproject.git
      version: master
      files:
        - src: subtool/test
          dst: ext/subtool.test/
          post_commands:
            - make

    - git: https://github.com/example/subtool2.git
      version: master
      dst: ext/subtool2/
      post_commands:
        - make

.. code-block:: bash

    $ gilt overlay

Display the git commands being executed.

.. code-block:: bash

    $ gilt --debug overlay

Use an alternate config file (default `gilt.yml`).

.. code-block:: bash

    $ gilt --config /path/to/gilt.yml overlay

Molecule
========

Integrates with `Molecule`_ as an `Ansible Galaxy`_ CLI replacement.  Update
`molecule.yml` with the following.

.. code-block:: yaml

    ---
    dependency:
      name: shell
      command: gilt overlay


Similar Tools
=============

* `Ansible Galaxy`_

.. toctree::
   :maxdepth: 3

   testing
   contributing
   development
   changelog
   authors
   autodoc

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`

.. _`Ansible Galaxy`: http://docs.ansible.com/ansible/galaxy.html#the-command-line-tool
.. _`Molecule`: http://molecule.readthedocs.io
