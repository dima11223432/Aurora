"""add_data_about_categories

Revision ID: e30013d5fb50
Revises: 77aa0dd86ab8
Create Date: 2026-05-09 16:15:10.323666

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "e30013d5fb50"
down_revision: Union[str, Sequence[str], None] = "77aa0dd86ab8"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.execute("INSERT INTO channel_categories (name) VALUES ('stocks'), ('crypto')")


def downgrade() -> None:
    op.execute("DELETE FROM channel_categories WHERE name IN ('stocks', 'crypto')")
