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

Testing
=======

.. code-block:: bash

  $ pip install tox
  $ tox

Similar Tools
=============

* `Reponimous`_
* `Ansible Galaxy`_

.. toctree::
   :maxdepth: 3

   autodoc

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`

.. _`Ansible Galaxy`: http://docs.ansible.com/ansible/galaxy.html#the-command-line-tool
.. _`Molecule`: http://molecule.readthedocs.io
.. _`Reponimous`: http://github.com/craigtracey/reponimous
