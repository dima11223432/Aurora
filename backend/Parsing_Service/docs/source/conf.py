import os
import sys
sys.path.insert(0, os.path.abspath('../../'))

project = 'Parsing_Service'
copyright = '2026, Ubenben'
author = 'Ubenben'
release = '1.0'

extensions = [
    'sphinx.ext.autodoc',
    'sphinx.ext.napoleon',
]

templates_path = ['_templates']
exclude_patterns = []

language = 'ru'

html_theme = 'alabaster'
html_static_path = ['_static']

autodoc_mock_imports = [
    "confluent_kafka",
    "dotenv",
    "loguru",
    "asyncpg",
    "psycopg2",
    "telethon",
    "socks",
]
