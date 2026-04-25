"""add_channels_categories

Revision ID: 599ed1df2acf
Revises: 98729658fa89
Create Date: 2026-04-25 17:08:30.732345

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "599ed1df2acf"
down_revision: Union[str, Sequence[str], None] = "98729658fa89"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "channel_categories",
        sa.Column("id", sa.Integer(), primary_key=True),
        sa.Column("name", sa.String(length=100), nullable=False, unique=True),
    )


def downgrade() -> None:
    op.drop_table("channel_categories")
