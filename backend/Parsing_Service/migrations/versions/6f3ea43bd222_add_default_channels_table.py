"""add_default_channels_table

Revision ID: 6f3ea43bd222
Revises: a735170311f8
Create Date: 2026-05-01 12:43:49.105529

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "6f3ea43bd222"
down_revision: Union[str, Sequence[str], None] = "a735170311f8"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "default_channels",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "channel_username",
            sa.String(255),
            nullable=False,
        ),
        sa.UniqueConstraint("channel_username", name="uq_default_channels"),
    )


def downgrade() -> None:
    op.drop_table("default_channels")
