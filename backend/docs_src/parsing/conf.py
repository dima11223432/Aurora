"""Sphinx configuration for Parsing Service."""

import os
import sys

sys.path.insert(0, os.path.abspath("/Users/dm_sfr/Aurora/backend/Parsing_Service"))

project = "Parsing Service"
copyright = "2026, Aurora"
author = "Aurora"
release = "1.0"

extensions = [
    "sphinx.ext.napoleon",
    "sphinx.ext.viewcode",
    "sphinx_autodoc_typehints",
]

templates_path = ["_templates"]
exclude_patterns = ["_build", "Thumbs.db", ".DS_Store"]

html_theme = "sphinx_rtd_theme"
html_static_path = ["_static"]

napoleon_google_docstring = True
napoleon_numpy_docstring = False
autodoc_typehints = "description"
